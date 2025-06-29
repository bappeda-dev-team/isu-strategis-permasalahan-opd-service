CREATE TABLE tb_permasalahan_opd(
    id INT PRIMARY KEY AUTO_INCREMENT,
    pokin_id INT NOT NULL,
    permasalahan VARCHAR(255) NOT NULL,
    level_pohon INT NOT NULL,
    jenis_masalah VARCHAR(255),
    kode_opd VARCHAR(255) NOT NULL,
    nama_opd VARCHAR(255) NOT NULL,
    tahun VARCHAR(255) NOT NULL,
    isu_strategis_id INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;