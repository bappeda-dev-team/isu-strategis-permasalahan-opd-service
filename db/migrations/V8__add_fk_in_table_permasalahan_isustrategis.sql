ALTER TABLE tb_permasalahan_isu_strategis
    ADD CONSTRAINT fk_permasalahan_isu_strategis_permasalahan_terpilih
    FOREIGN KEY (id_permasalahan) 
    REFERENCES tb_permasalahan_terpilih(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE tb_permasalahan_isu_strategis
    ADD CONSTRAINT fk_permasalahan_isu_strategis_isu_strategis
    FOREIGN KEY (id_isu_strategis) 
    REFERENCES tb_isu_strategis_opd(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;