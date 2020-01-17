using DotNetEnv;
using Microsoft.Extensions.Logging;

namespace RedirectProtect.Database
{
    public class RedirectProtectConfig : IRedirectProtectConfig
    {
        public string CollectionName { get; set; }
        public string DatabaseName { get; set; }
        public string ConnectionString { get; set; }
    }
    public interface IRedirectProtectConfig
    {
        string CollectionName { get; set; }
        string ConnectionString { get; set; }
        string DatabaseName { get; set; }
    }
}