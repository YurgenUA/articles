using Amazon.Lambda.Core;
using System.Collections;
using System;
using System.Text;
using System.Net.Http;
using DewCore.RestClient;
using DewCore.Abstract.Internet;
using System.Threading.Tasks;
using Newtonsoft.Json;

[assembly: LambdaSerializer(typeof(Amazon.Lambda.Serialization.Json.JsonSerializer))]
namespace AwsDotnetCsharp
{
  public class Handler
  {
    // entry point for rotation Lambda
    public async Task<bool> Handle(object request, ILambdaContext context)
    {
      try
      {
        var tokenGenerator = new TokenGenerator();
        // 1. get fresh token
        string newToken = await tokenGenerator.Generate();
        context.Logger.Log(newToken);
        if (!string.IsNullOrWhiteSpace(newToken))
        {
        context.Logger.Log("Storing token...");
          // 2. store fresh token into AWS Secrets Service.
          return await tokenGenerator.Store(newToken);
        }
        return false;
      }
      catch (Exception e)
      {
        context.Logger.Log($"Exception occurred :{e.Message}");
        throw new ApplicationException("Lambda execution exception", e);
      }
    }
  }

  public class TokenGenerator
  {
    // POCO class used to parse response from Auth0 API 
    class Auth0Response
    {
      [JsonProperty("access_token")]
      public string accessToken { get; set; }
      [JsonProperty("expires_in")]
      public string expiresIn { get; set; }
      [JsonProperty("token_type")]
      public string tokenType { get; set; }
    }

    // POCO class used to parse response from AWS Secrets Manager with Auth0 credentials 
    class Auth0SecretsResponse
    {
      [JsonProperty("CLIENT_ID")]
      public string clientId { get; set; }
      [JsonProperty("CLIENT_SECRET")]
      public string clientSecret { get; set; }
    }

    class JWTSecretsRequest
    // POCO class used to send fresh token to AWS Secrets Manager 
    {
      [JsonProperty("GENERATED_AT")]
      public string generatedAt { get; set; }
      [JsonProperty("BEARER_TOKEN")]
      public string bearerToken { get; set; }
    }

    const string AUTH0_SECRET_NAME = "auth0-secrets-demo";
    const string TOKEN_SECRET_NAME = "service-token-demo";
    private static string Auth0URL
    {
      get
      {
        return System.Environment.GetEnvironmentVariable("AUTH0_URL");
      }
    }
    private static string Auth0Audience
    {
      get
      {
        return System.Environment.GetEnvironmentVariable("AUTH0_AUDIENCE");
      }
    }

    public async Task<string> Generate()
    {
      // get Auth0 secrets from AWS SecretsManager entry
      var auth0Secrets = await SecretsManagerWrapper.GetSecret<Auth0SecretsResponse>(AUTH0_SECRET_NAME);

      // prepare Auth0 JWT service HTTP call 
      RESTRequest request = new RESTRequest(TokenGenerator.Auth0URL);
      request.SetMethod(Method.POST);
      var bodyString = JsonConvert.SerializeObject(new
      {
        client_id = auth0Secrets.clientId,
        client_secret = auth0Secrets.clientSecret,
        audience = TokenGenerator.Auth0Audience,
        grant_type = "client_credentials"
      });
      request.AddContent(new StringContent(bodyString, Encoding.UTF8, "application/json"));

      // make HTTP call
      using (RESTClient client = new RESTClient())
      {
        using (RESTResponse response = (RESTResponse)await client.PerformRequestAsync(request))
        {
          if (response.IsSuccesStatusCode())
          {
            string responseString = await response.ReadResponseAsStringAsync();
            Console.WriteLine(responseString);
            var responseObject = JsonConvert.DeserializeObject<Auth0Response>(responseString);
            return responseObject.accessToken;
          }
          throw new ApplicationException("Failed to get JWT from Auth0");
        }
      }
    }

    public async Task<bool> Store(string jwt)
    {
      var content = new JWTSecretsRequest
      {
        generatedAt = DateTime.UtcNow.ToString(),
        bearerToken = $"Bearer {jwt}"
      };
      return await SecretsManagerWrapper.SetSecret<JWTSecretsRequest>(TOKEN_SECRET_NAME, content);
    }
  }
}
