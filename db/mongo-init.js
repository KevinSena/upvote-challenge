conn = new Mongo();
db = conn.getDB("upvote");

db.posts.drop();
db.posts.insertMany([
  { "title": "I love java", "desc": "it's fake", "votes": 3, "author": "Fulana Pereira" },
  { "title": "I'm anxious to the Rings of Power", "desc": "I expect it to be cool", "votes": 5, "author": "Myself" },
  { "title": "Have you ever seen the rain?", "desc": "Old but good", "votes": 0, "author": "Bertrano" },
  { "title": "Hello World!", "desc": "good start", "votes": 10, "author": "Ciclano" },
  { "title": "Mistborn", "desc": "Would be a good game", "votes": 6, "author": "Fulano New" },
]);
