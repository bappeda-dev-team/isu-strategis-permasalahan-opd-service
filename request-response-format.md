# Request Response format

## Permasalahan

#### [REQUEST-BODY] JSON body format for POST and PATCH

```json
{
    "id_pokin": "123",
    "pokin": "Penguatan akuntabilitas dalam hal perencanaan penganggaran"
    "kode_opd": "123.456",
    "tahun": "2025",
    "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
    "keterangan": "ambil dari pokin 2025",
    "jenis_masalah": "MasalahPokok",
}
```

#### [RESPONSE-BODY] success input

```json
{
    "id": 1,
    "id_pokin": "123",
    "pokin": "Penguatan akuntabilitas dalam hal perencanaan penganggaran"
    "kode_opd": "123.456",
    "tahun": "2025",
    "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
    "keterangan": "ambil dari pokin 2025",
    "jenis_masalah": "MasalahPokok",
    "is_terpilih": false,
    "created_date": "2025-06-03 02:43:50",
    "last_modified_date": "2025-06-03 02:43:50",
}
```

#### [RESPONSE-BODY] get permasalahan by id

```json
{
    "id": 1,
    "id_pokin": "123",
    "pokin": "Penguatan akuntabilitas dalam hal perencanaan penganggaran"
    "kode_opd": "123.456",
    "tahun": "2025",
    "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
    "keterangan": "ambil dari pokin 2025",
    "jenis_masalah": "MasalahPokok",
    "is_terpilih": false,
    "created_date": "2025-06-03 02:43:50",
    "last_modified_date": "2025-06-03 02:43:50",
}
```

#### [RESPONSE-BODY] get permasalahan by opd dan tahun

```json
{
    "kode_opd": "123.456",
    "tahun": "2025",
    "permasalahan_opd": [
        {
        "id": 1,
        "id_pokin": "123",
        "pokin": "Penguatan akuntabilitas dalam hal perencanaan penganggaran"
        "kode_opd": "123.456",
        "tahun": "2025",
        "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
        "keterangan": "ambil dari pokin 2025",
        "jenis_masalah": "MasalahPokok",
        "is_terpilih": false,
        "created_date": "2025-06-03 02:43:50",
        "last_modified_date": "2025-06-03 02:43:50",
        },
        {
        "id": 11,
        "id_pokin": "1231",
        "pokin": "penganggaran"
        "kode_opd": "123.456",
        "tahun": "2025",
        "masalah": "penganggaran masih belum optimal"
        "keterangan": "ambil dari pokin 2025",
        "jenis_masalah": "Masalah",
        "is_terpilih": false,
        "created_date": "2025-06-04 02:43:50",
        "last_modified_date": "2025-06-04 02:43:50",
        },
    ]
}
```

#### [RESPONSE-BODY] permasalahan dipilih menjadi permasalahan terpilih

```json
{
    "id": 1,
    "id_pokin": "123",
    "pokin": "Penguatan akuntabilitas dalam hal perencanaan penganggaran"
    "kode_opd": "123.456",
    "tahun": "2025",
    "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
    "keterangan": "ambil dari pokin 2025",
    "jenis_masalah": "MasalahPokok",
    "is_terpilih": true,
    "created_date": "2025-06-03 02:43:50",
    "last_modified_date": "2025-06-03 02:43:50",
}
```

---

## Permasalahan Terpilih

#### [RESPONSE-BODY] get permasalahan terpilih by kode_opd dan tahun

```json
{
  "kode_opd": "123.456",
  "tahun": "2025",
  "permasalahan_terpilih": [
    {
      "id": 9,
      "id_permasalahan_opd": 1,
      "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
    },
    {
      "id": 99,
      "id_permasalahan_opd": 11,
      "masalah": "penganggaran masih belum optimal"
    }
  ]
}
```

---

## Isu Strategis

#### [REQUEST-BODY] POST and PATCH

```json
{
  "kode_opd": "123.456",
  "tahun": "2025",
  "kode_bidang_urusan": "123",
  "nama_bidang_urusan": "BIDANG URUSAN CONTOH A",
  "isu_strategis": "Peningkatan Kualitas Perencanaan Inovasi Daerah",
  "permasalahans": [
    {
      "id_permasalahan": 11,
      "data_dukungs": [
        {
          "nama_data_dukung": "rekap anggaran opd",
          "narasi_data_dukung": "rekapitulasi realisasi capaian anggaran opd",
          "jumlah_data": [
            {
              "tahun": "2021",
              "jumlah": "33",
              "satuan": "dokumen"
            },
            {
              "tahun": "2022",
              "jumlah": "33",
              "satuan": "dokumen"
            },
            {
              "tahun": "2023",
              "jumlah": "33",
              "satuan": "dokumen"
            }
          ]
        }
      ]
    }
  ]
}
```

#### [RESPONSE-BODY] success response

```json
{
    "id": 3,
    "kode_opd": "123.456",
    "tahun": "2025",
    "kode_bidang_urusan": "123",
    "nama_bidang_urusan": "BIDANG URUSAN CONTOH A",
    "isu_strategis": "Peningkatan Kualitas Perencanaan Inovasi Daerah",
    "created_date": "2025-06-03 02:43:50",
    "last_modified_date": "2025-06-03 02:43:50",
    "permasalahans": [
    {
        "id": 9,
        "id_permasalahan_opd": 1,
        "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
        "data_dukungs": [
        {
            "id": 100,
            "id_permasalahan_opd": 1,
            "nama_data_dukung": "rekapitulasi penyusunan proposal inovasi",
            "narasi_data_dukung": "Proposal inovasi yang berasal dari masyarakat....",
            "jumlah_data": [
            {
                "id": 7,
                "id_data_dukung": 100,
                "tahun": "2021",
                "jumlah": "100",
                "satuan": "inovasi"
            },
            {
                "id": 8,
                "id_data_dukung": 100,
                "tahun": "2022",
                "jumlah": "101",
                "satuan": "inovasi"
            },
            {
                "id": 9,
                "id_data_dukung": 100,
                "tahun": "2022",
                "jumlah": "103",
                "satuan": "inovasi"
            },
            {
                "id": 10,
                "id_data_dukung": 100,
                "tahun": "2023",
                "jumlah": "104",
                "satuan": "inovasi"
            }
            ]
        },
        {
            "id": 101,
            "id_permasalahan_opd": 1,
            "nama_data_dukung": "reakpitulasi proposal inovasi asn",
            "narasi_data_dukung": "Proposal inovasi dari internal pemkot dan dilakukan oleh asn atau mitra....",
            "jumlah_data": [
            {
                "id": 77,
                "id_data_dukung": 101,
                "tahun": "2021",
                "jumlah": "110",
                "satuan": "inovasi"
            },
            {
                "id": 78,
                "id_data_dukung": 101,
                "tahun": "2022",
                "jumlah": "111",
                "satuan": "inovasi"
            },
            {
                "id": 79,
                "id_data_dukung": 101,
                "tahun": "2022",
                "jumlah": "113",
                "satuan": "inovasi"
            },
            {
                "id": 80,
                "id_data_dukung": 101,
                "tahun": "2023",
                "jumlah": "114",
                "satuan": "inovasi"
            }
            ]
        }
        ]
    },
    {
        "id": 99,
        "id_permasalahan_opd": 11,
        "masalah": "penganggaran masih belum optimal"
        "data_dukungs": [
        {
            "id": 30,
            "id_permasalahan_opd": 11,
            "nama_data_dukung": "rekap anggaran opd",
            "narasi_data_dukung": "rekapitulasi realisasi capaian anggaran opd",
            "jumlah_data": [
            {
                "id": 55,
                "id_data_dukung": 30,
                "tahun": "2021",
                "jumlah": "33",
                "satuan": "dokumen"
            },
            {
                "id": 56,
                "id_data_dukung": 30,
                "tahun": "2022",
                "jumlah": "33",
                "satuan": "dokumen"
            },
            {
                "id": 57,
                "id_data_dukung": 30,
                "tahun": "2023",
                "jumlah": "33",
                "satuan": "dokumen"
            },
            ]
        }
        ]
    }
    ]
}
```

#### [RESPONSE-BODY] get all by kode_opd dan tahun

```json
{
    "kode_opd": "123.456",
    "tahun": "2025",
    "isu_strategis_opds": [
        {
        "id": 3,
        "kode_opd": "123.456",
        "tahun": "2025",
        "kode_bidang_urusan": "123",
        "nama_bidang_urusan": "BIDANG URUSAN CONTOH A",
        "isu_strategis": "Peningkatan Kualitas Perencanaan Inovasi Daerah",
        "created_date": "2025-06-03 02:43:50",
        "last_modified_date": "2025-06-03 02:43:50",
        "permasalahans": [
            {
            "id": 9,
            "id_permasalahan_opd": 1,
            "masalah": "Akuntabilitas dalam hal perencanaan penganggaran masih belum optimal"
            "data_dukungs": [
                {
                "id": 100,
                "id_permasalahan_opd": 1,
                "nama_data_dukung": "rekapitulasi penyusunan proposal inovasi",
                "narasi_data_dukung": "Proposal inovasi yang berasal dari masyarakat....",
                "jumlah_data": [
                    {
                    "id": 7,
                    "id_data_dukung": 100,
                    "tahun": "2021",
                    "jumlah": "100",
                    "satuan": "inovasi"
                    },
                    {
                    "id": 8,
                    "id_data_dukung": 100,
                    "tahun": "2022",
                    "jumlah": "101",
                    "satuan": "inovasi"
                    },
                    {
                    "id": 9,
                    "id_data_dukung": 100,
                    "tahun": "2022",
                    "jumlah": "103",
                    "satuan": "inovasi"
                    },
                    {
                    "id": 10,
                    "id_data_dukung": 100,
                    "tahun": "2023",
                    "jumlah": "104",
                    "satuan": "inovasi"
                    }
                ]
                },
                {
                "id": 101,
                "id_permasalahan_opd": 1,
                "nama_data_dukung": "reakpitulasi proposal inovasi asn",
                "narasi_data_dukung": "Proposal inovasi dari internal pemkot dan dilakukan oleh asn atau mitra....",
                "jumlah_data": [
                    {
                    "id": 77,
                    "id_data_dukung": 101,
                    "tahun": "2021",
                    "jumlah": "110",
                    "satuan": "inovasi"
                    },
                    {
                    "id": 78,
                    "id_data_dukung": 101,
                    "tahun": "2022",
                    "jumlah": "111",
                    "satuan": "inovasi"
                    },
                    {
                    "id": 79,
                    "id_data_dukung": 101,
                    "tahun": "2022",
                    "jumlah": "113",
                    "satuan": "inovasi"
                    },
                    {
                    "id": 80,
                    "id_data_dukung": 101,
                    "tahun": "2023",
                    "jumlah": "114",
                    "satuan": "inovasi"
                    }
                ]
                }
            ]
            },
            {
            "id": 99,
            "id_permasalahan_opd": 11,
            "masalah": "penganggaran masih belum optimal"
            "data_dukungs": [
                {
                "id": 30,
                "id_permasalahan_opd": 11,
                "nama_data_dukung": "rekap anggaran opd",
                "narasi_data_dukung": "rekapitulasi realisasi capaian anggaran opd",
                "jumlah_data": [
                    {
                    "id": 55,
                    "id_data_dukung": 30,
                    "tahun": "2021",
                    "jumlah": "33",
                    "satuan": "dokumen"
                    },
                    {
                    "id": 56,
                    "id_data_dukung": 30,
                    "tahun": "2022",
                    "jumlah": "33",
                    "satuan": "dokumen"
                    },
                    {
                    "id": 57,
                    "id_data_dukung": 30,
                    "tahun": "2023",
                    "jumlah": "33",
                    "satuan": "dokumen"
                    },
                ]
                }
            ]
            }
        ]
        }
    ]
}
```

---
