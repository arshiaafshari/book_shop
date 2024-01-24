package query

var GetBookQuery = " SELECT *"+ 
				    "FROM (  SELECT book_id, Min(s_price) AS price"+ 
	  					 	"FROM quantity_book "+
	   						"WHERE quantity > 0"+
	  						"GROUP BY id) AS t1"+
				    "RIGHT JOIN (	 SELECT id, title, author, publisher, year, language, isbc, edition, pages, translator"+ 
									"FROM books"+ 
									"where deleted_at is null ) AS t2"+
					"ON t1.id = t2.id"





//"SELECT * FROM books LEFT JOIN quntity_book USING(book_id) WHERE books.deleted_at is null AND quntity_book.quantity > 0"					