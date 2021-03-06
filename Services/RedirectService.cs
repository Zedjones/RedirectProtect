using MongoDB.Driver;
using RedirectProtect.Database;
using System.Collections.Generic;
using Microsoft.Extensions.Logging;
using RedirectProtect.Database.Models;
using Newtonsoft.Json;
using System;

namespace RedirectProtect.Services
{
    public interface IRedirectService
    {
        bool RedirectExists(Redirect redir);
        Redirect GetRedirect(string path);
        List<Redirect> GetRedirects();
        void DeleteRedirect(string path);
        void DeleteRedirect(Redirect redirIn);
        Redirect Create(RedirectDto redirect);
    }
    public class RedirectService : IRedirectService
    {
        private readonly IMongoCollection<Redirect> _redirects;
        private readonly ILogger<RedirectService> _logger;
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
            _logger = logger;
        }
        private bool PathExists(string path) =>
            _redirects.Find(redirect => redirect.Path == path).CountDocuments() == 1;

        public bool RedirectExists(Redirect redir) =>
            _redirects.Find(redirect => redirect.Id == redir.Id && redirect.Path == redir.Path).CountDocuments() == 1;

        public Redirect GetRedirect(string path) =>
            _redirects.Find(redirect => redirect.Path == path).FirstOrDefault();

        public List<Redirect> GetRedirects() =>
            _redirects.Find(_ => true).ToList();

        public void DeleteRedirect(string path) =>
            _redirects.DeleteOne(redirect => redirect.Path == path);

        public void DeleteRedirect(Redirect redirIn) =>
            _redirects.DeleteOne(redir => redir.Path == redirIn.Path);

        public Redirect Create(RedirectDto redirect)
        {
            String path;
            do
            {
                path = Utils.RandomHex.GetRandomHexNumber(8).ToLower();
            }
            while (PathExists(path));

            var hashedPass = BCrypt.Net.BCrypt.HashPassword(redirect.Password);
            DateTime? fullTime;
            if (redirect.TTL is null)
            {
                fullTime = null;
            }
            else
            {
                fullTime = DateTime.UtcNow.AddHours(redirect.TTL.Value.Hour).AddMinutes(redirect.TTL.Value.Minute);
            }

            var newRedir = new Redirect
            {
                ExpirationTime = fullTime,
                Password = hashedPass,
                URL = redirect.URL,
                Path = path
            };
            _redirects.InsertOne(newRedir);
            _logger.LogInformation($"Inserted redirect: {JsonConvert.SerializeObject(newRedir)}");
            return newRedir;
        }
    }
}