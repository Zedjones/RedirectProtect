using MongoDB.Driver;
using RedirectProtect.Database;
using System.Collections.Generic;
using Microsoft.Extensions.Logging;
using RedirectProtect.Database.Models;
using System;

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
            else if (settings.ConnectionString is null)
            {
                logger.LogCritical("Missing connection string in configuration");
                somethingMissing = true;
            }
            else if (settings.DatabaseName is null)
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

        public void Create(RedirectDto redirect)
        {
            _redirects.InsertOne(new Redirect
            {
                TTL = redirect.TTL,
                Password = redirect.Password,
                CreatedOn = DateTime.Now,
                URL = redirect.URL
            });
        }
    }
}