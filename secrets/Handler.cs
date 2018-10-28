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
    public async Task<bool> Handle(object request, ILambdaContext context)
    {
      try
      {
        Console.WriteLine($"env - {System.Environment.GetEnvironmentVariable("CLIENT_ID")}");
        string newToken = await new TokenGenerator().Generate();
        context.Logger.Log(newToken);
        return true;
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
    class Auth0SecretsResponse
    {
      [JsonProperty("CLIENT_ID")]
      string clientId { get; set; }
      [JsonProperty("CLIENT_SECRET")]
      string clientSecret { get; set; }
    }

    const string AUTH0_SECRET_NAME = "auth0-secrets-demo";
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
        client_id = TokenGenerator.auth0Secrets.clientId,
        client_secret = TokenGenerator.auth0Secrets.clientSecret,
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
            return responseString;
          }
          throw new ApplicationException("Failed to get JWT from Auth0");
        }
      }
    }
  }
}
