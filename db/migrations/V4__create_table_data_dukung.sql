CREATE TABLE tb_data_dukung (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_permasalahan INT NOT NULL,
    nama_data_dukung TEXT,
    narasi_data_dukung TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;