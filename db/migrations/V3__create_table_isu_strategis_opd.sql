CREATE TABLE tb_isu_strategis_opd (
    id INT PRIMARY KEY AUTO_INCREMENT,
    kode_opd VARCHAR(255) NOT NULL,
    nama_opd VARCHAR(255) NOT NULL,
    kode_bidang_urusan VARCHAR(255) NOT NULL,
    nama_bidang_urusan VARCHAR(255) NOT NULL,
    tahun_awal VARCHAR(255) NOT NULL,
    tahun_akhir VARCHAR(255) NOT NULL,
    isu_strategis TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;