using System;

using MongoDB.Bson;

namespace RedirectProtect.Database.Models
{
    public class Redirect
    {
        public ObjectId Id { get; set; }
        public DateTime CreatedOn { get; set; }
        public string Path { get; set; }
        public string URL { get; set; }
        public string Password { get; set; }
        public DateTime TTL { get; set; }
    }
}