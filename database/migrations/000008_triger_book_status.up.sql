CREATE TRIGGER book_status_on_loan
AFTER INSERT ON loans
FOR EACH ROW
EXECUTE FUNCTION update_book_status_on_loan();
