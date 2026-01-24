ALTER TABLE tb_data_dukung
ADD COLUMN id_isustrategis INT NULL,
ADD CONSTRAINT fk_data_dukung_isustrategis
FOREIGN KEY (id_isustrategis)
REFERENCES tb_isu_strategis_opd(id)
ON DELETE CASCADE
ON UPDATE CASCADE;