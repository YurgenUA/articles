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

    public async static Task<T> GetSecret<T>(string secretName) where T : new()
    {
      MemoryStream memoryStream = new MemoryStream();

      IAmazonSecretsManager client = new AmazonSecretsManagerClient(RegionEndpoint.GetBySystemName(SecretsManagerWrapper.region));

      GetSecretValueRequest request = new GetSecretValueRequest();
      request.SecretId = secretName;
      request.VersionStage = "AWSCURRENT";

      GetSecretValueResponse response = null;
      string secret = String.Empty;

      response = await client.GetSecretValueAsync(request);
      if (response.SecretString != null)
      {
        secret = response.SecretString;
        return JsonConvert.DeserializeObject<T>(secret);
      }
      throw new ApplicationException($"Failed type in AWS Secret {secretName}");
    }

  }
}