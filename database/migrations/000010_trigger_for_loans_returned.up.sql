CREATE TRIGGER book_status_on_return
AFTER UPDATE OF return_date ON loans
FOR EACH ROW
WHEN (NEW.return_date IS NOT NULL)
EXECUTE FUNCTION update_book_status_on_return();
