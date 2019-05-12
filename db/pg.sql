CREATE TABLE article (
  id serial NOT NULL PRIMARY KEY,
  title varchar(200) NOT NULL,
  body text,
  created_on date NOT NULL DEFAULT NOW()
);


CREATE TABLE tag (
  id serial NOT NULL PRIMARY KEY,
  article_id int NOT NULL REFERENCES article(id),
  created_on date NOT NULL DEFAULT NOW(),
  tag_name varchar(50) NOT NULL,
  related_tags varchar[]
);