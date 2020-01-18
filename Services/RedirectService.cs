using MongoDB.Driver;
using RedirectProtect.Database;
using System.Collections.Generic;
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