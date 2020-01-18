using System;

using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace RedirectProtect.Database.Models
{
    public class Redirect
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string Id { get; set; }
        public DateTime CreatedOn { get; set; }
        public string Path { get; set; }
        public string URL { get; set; }
        public string Password { get; set; }
        public string TTL { get; set; }
    }
}