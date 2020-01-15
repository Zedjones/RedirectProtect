using DotNetEnv;
using Microsoft.Extensions.Logging;

namespace RedirectProtect.Database
{
    public class Config
    {
        public readonly string Username;
        public readonly string Password;
        public readonly string ConnectionString;
        private readonly ILogger _logger;
        public Config(ILogger<Config> logger)
        {
            _logger = logger;
            Env.Load(new DotNetEnv.Env.LoadOptions(
                parseVariables: true
            ));
            Username = Env.GetString("MONGO_INITDB_ROOT_USERNAME");
            Password = Env.GetString("MONGO_INITDB_ROOT_PASSWORD");
            ConnectionString = Env.GetString("CONNECTION_STRING");
            _logger.LogInformation("test");
        }
    }
}