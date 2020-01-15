using MongoDB.Driver;

namespace RedirectProtect.Database
{
    public class MongoClient
    {
        private MongoDB.Driver.MongoClient client;
        public MongoClient()
        {
            client = new MongoDB.Driver.MongoClient(Config.ConnectionString);
        }
    }
}