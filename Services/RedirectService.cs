using MongoDB.Driver;
using RedirectProtect.Database;
using RedirectProtect.Database.Models;

namespace RedirectProtect.Services
{
    public class RedirectService
    {
        private readonly IMongoCollection<Redirect> _redirects;
        public RedirectService(IRedirectProtectConfig settings) 
        {
            var client = new MongoClient(settings.ConnectionString);
            var database = client.GetDatabase(settings.DatabaseName);

            _redirects = database.GetCollection<Redirect>(settings.CollectionName);
        }
    }
}