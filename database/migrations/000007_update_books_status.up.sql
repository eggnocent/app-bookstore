CREATE OR REPLACE FUNCTION update_book_status_on_loan() 
RETURNS TRIGGER AS $$
BEGIN
    -- Log untuk debugging (bisa dihapus jika tidak diperlukan)
    RAISE NOTICE 'Trigger executed: Book % set to borrowed', NEW.book_id;

    -- Update status buku menjadi 'borrowed'
    UPDATE books 
    SET status = 'borrowed' 
    WHERE id = NEW.book_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
