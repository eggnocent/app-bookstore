CREATE UNIQUE INDEX unique_loan_book ON loans (book_id) WHERE return_date IS NULL;
