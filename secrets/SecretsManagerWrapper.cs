using System;
using System.IO;
using System.Threading.Tasks;
using Amazon;
using Amazon.SecretsManager;
using Amazon.SecretsManager.Model;
using Newtonsoft.Json;


namespace AwsDotnetCsharp
{

  public class SecretsManagerWrapper
  {
    private static readonly string region = "eu-west-1";

    //////////////////
    // Return secret object from AWS Secret
    //////////////////
    public async static Task<T> GetSecret<T>(string secretName)
    {
      MemoryStream memoryStream = new MemoryStream();

      // Accessing AWS Secrets Manager
      IAmazonSecretsManager client = new AmazonSecretsManagerClient(RegionEndpoint.GetBySystemName(SecretsManagerWrapper.region));

      GetSecretValueRequest request = new GetSecretValueRequest();
      request.SecretId = secretName;
      request.VersionStage = "AWSCURRENT";

      GetSecretValueResponse response = null;
      string secret = String.Empty;

      // make request
      response = await client.GetSecretValueAsync(request);
      if (response.SecretString != null)
      {
        secret = response.SecretString;
        // expecting secret to be saved as JSON
        return JsonConvert.DeserializeObject<T>(secret);
      }
      throw new ApplicationException($"Failed type in AWS Secret {secretName}");
    }

    //////////////////
    // Store secret object to AWS Secret
    //////////////////
     public async static Task<bool> SetSecret<T>(string secretName, T secretObject)
    {
      MemoryStream memoryStream = new MemoryStream();
      // Accessing AWS Secrets Manager
      IAmazonSecretsManager client = new AmazonSecretsManagerClient(RegionEndpoint.GetBySystemName(SecretsManagerWrapper.region));

      // store sectret as JSON string
      var sectedAsJson = JsonConvert.SerializeObject(secretObject);

      PutSecretValueRequest request = new PutSecretValueRequest();
      request.SecretId = secretName;
      request.SecretString = sectedAsJson;

      PutSecretValueResponse response = null;

      // make store request
      response = await client.PutSecretValueAsync(request);
      return response.ContentLength > 0;
    }
 }
}