-- 1. Relasi untuk tb_permasalahan_terpilih
ALTER TABLE tb_permasalahan_terpilih
    ADD CONSTRAINT fk_permasalahan_terpilih_permasalahan_opd_v6
    FOREIGN KEY (permasalahan_opd_id) 
    REFERENCES tb_permasalahan_opd(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- 2. Relasi untuk tb_data_dukung
ALTER TABLE tb_data_dukung
    ADD CONSTRAINT fk_data_dukung_permasalahan_v6
    FOREIGN KEY (id_permasalahan) 
    REFERENCES tb_permasalahan_opd(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- 3. Relasi untuk tb_jumlah_data
ALTER TABLE tb_jumlah_data
    ADD CONSTRAINT fk_jumlah_data_data_dukung_v6
    FOREIGN KEY (id_data_dukung) 
    REFERENCES tb_data_dukung(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- 4. Tambah kolom dan foreign key untuk isu_strategis di tb_permasalahan_opd
ALTER TABLE tb_permasalahan_opd
    MODIFY COLUMN isu_strategis_id INT DEFAULT 0;
    -- Hapus foreign key constraint karena kita akan menggunakan trigger saja

-- 5. Trigger untuk menghapus isu strategis jika permasalahan terkait dihapus
DELIMITER //

CREATE TRIGGER after_permasalahan_delete 
AFTER DELETE ON tb_permasalahan_opd
FOR EACH ROW
BEGIN
    -- Jika permasalahan yang dihapus memiliki isu strategis (tidak 0)
    IF OLD.isu_strategis_id != 0 THEN
        -- Cek apakah masih ada permasalahan lain dengan isu strategis yang sama
        IF NOT EXISTS (
            SELECT 1 
            FROM tb_permasalahan_opd 
            WHERE isu_strategis_id = OLD.isu_strategis_id
        ) THEN
            -- Jika tidak ada, hapus isu strategis
            DELETE FROM tb_isu_strategis_opd 
            WHERE id = OLD.isu_strategis_id;
        END IF;
    END IF;
END//

DELIMITER ;