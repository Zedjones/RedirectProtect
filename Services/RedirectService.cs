using MongoDB.Driver;
using RedirectProtect.Database;
using System.Collections.Generic;
using Microsoft.Extensions.Logging;
using RedirectProtect.Database.Models;
using Microsoft.Extensions.Hosting;

namespace RedirectProtect.Services
{
    public class RedirectService
    {
        private readonly IMongoCollection<Redirect> _redirects;
        public RedirectService(IRedirectProtectConfig settings, ILogger<RedirectService> logger)
        {
            var somethingMissing = false;
            if (settings.CollectionName is null)
            {
                logger.LogCritical("Missing collection name in configuration");
                somethingMissing = true;
            }
            else if(settings.ConnectionString is null)
            {
                logger.LogCritical("Missing connection string in configuration");
                somethingMissing = true;
            }
            else if(settings.DatabaseName is null)
            {
                logger.LogCritical("Missing database name in configuration");
                somethingMissing = true;
            }
            if (somethingMissing)
            {
                throw new MongoConfigurationException("Missing some values");
            }

            var client = new MongoClient(settings.ConnectionString);
            var database = client.GetDatabase(settings.DatabaseName);
            logger.LogInformation("Created Mongo Client");

            _redirects = database.GetCollection<Redirect>(settings.CollectionName);
        }
        public List<Redirect> GetRedirects() =>
            _redirects.Find(_ => true).ToList();

        public Redirect GetRedirect(string id) =>
            _redirects.Find<Redirect>(redirect => redirect.Id == id).FirstOrDefault();

        public void Create(Redirect redirect) =>
            _redirects.InsertOne(redirect);

        public void Remove(Redirect redirect) =>
            _redirects.DeleteOne(redir => redir.Id == redirect.Id);

        public void Remove(string id) =>
            _redirects.DeleteOne(redirect => redirect.Id == id);

    }
}