CREATE TABLE administrater
(
		username Char VARYING(20) NOT NULL,
		passwd  Char VARYING(20) NOT NULL,

		PRIMARY Key(username),
		UNIQUE(username),
		check(length(passwd)>8 and passwd similar to '%[A-Z]%'
		and passwd similar to '%[a-z]%' and passwd similar to '%[0-9]%')
);


CREATE TABLE books
(
		bookname Char VARYING(100) NOT NULL,
		bookid Char VARYING(100) NOT NULL,
		numberleft Integer,

		PRIMARY KEY(bookid),
		UNIQUE(bookid),
		check(numberleft >= 0)
);

INSERT INTO administrater VALUES('admin','A2x123456');


INSERT INTO books VALUES('Ansi Common Lisp','isbn0001',3);
INSERT INTO books VALUES('Clojure编程','isbn0002',5);
INSERT INTO books VALUES('SICP','isbn0003',10);
INSERT INTO books VALUES('On Lisp','isbn0004',2);
INSERT INTO books VALUES('On Lisp','isbn0004',2);

CREATE TABLE users
(
		username Char VARYING(20) NOT NULL,
		passwd  Char VARYING(20) NOT NULL,
		num integer,
        book		Char VARYING(300),
		bookname		Char VARYING(300),
		PRIMARY Key(username),
		UNIQUE(username),
		check(length(passwd)>8)
);


--------------------------------------------------------
SELECT * FROM users;
SELECT * FROM books ;
DELETE  FROM users WHERE username='bevis';
INSERT INTO users VALUES ('1','2',0,'');


UPDATE users set num=0,book='' where username='bevis';

UPDATE books set numberleft=numberleft-1 where bookid='isbn0003';
