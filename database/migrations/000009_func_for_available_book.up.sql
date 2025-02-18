CREATE OR REPLACE FUNCTION update_book_status_on_return() 
RETURNS TRIGGER AS $$
BEGIN
    -- Log untuk debugging
    RAISE NOTICE 'Trigger executed: Book % set to available', NEW.book_id;

    -- Update status buku menjadi 'available' jika buku dikembalikan
    UPDATE books 
    SET status = 'available' 
    WHERE id = NEW.book_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
