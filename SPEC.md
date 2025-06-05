# Spec isu dan permasalahan opd

- [ ] **Define Feature Summary**

  - [ ] Feature name: Isu-Permasalahan-Opd-Service
  - [ ] Service name: Isu-Permasalahan-Opd-Service
  - [ ] description: Backend service untuk isu dan permaslahaan opd
  - [ ] Goal or outcome:
    - User admin bisa input data permasalahan dari pohon kinerja opd
    - User admin dapat menetapakan `permasalahan terpilih`
    - User admin dapat input isu strategis opd
    - User admin dapat menambahkan `permasalahan terpilih` kedalam isu strategis opd

- [ ] **Write Context / Background**

  - [ ] Why the feature is needed: Untuk memenuhi kebutuhan isu strategis opd berdasarkan permasalahan yang ditemukan dalam opd.
  - [ ] What user or business problem it solves: Data untuk rpjmd, renstra, renja, serta bahan evaluasi kinerja opd
  - [ ] Any related features or legacy concerns: Pohon Kinerja OPD

- [ ] **List Business Requirements**

  - [ ] Expected user actions (what can user do?):
    - Create permasalahan dari pohon kinerja
    - Edit dan Update permasalahan
    - Pilih permasalahan
    - Create isu strategis opd
    - Edit dan Update isu strategis opd
    - Hapus isu strategis opd
    - Create data_dukung permasalahan opd
    - Update data_dukung permasalahan opd
    - Hapus data_dukung permasalahan opd
  - [ ] Expected results or outputs:
    - Permasalahan terpilih opd dalam tahun x
    - Isu strategis opd opd dalama tahun x
    - Data dukung permasalahan dalam tahun x (multiple tahun atau tahun dalam satu periode)
  - [ ] Must-have filters, validations, or restrictions:
    - [ ] filter: filter permasalahan by opd dan tahun untuk admin opd.
    - [ ] filter: filter isu strategis by opd dan tahun untuk admin opd.
    - [ ] validasi: semua field harus terisi kecuali id_pokin dan pokin, jika tidak status -> belum lengkap.
    - [ ] restrictions: permasalahan hanya dapat dipilih menjadi permasalahan terpilih satu kali.
    - [ ] restrictions: permasalahan terpilih hanya dapat digunakan dalam satu isu strategis.
  - [ ] Logging or audit requirements: -

- [ ] **Define API Contract**

  - [ ] HTTP method and endpoint path
        Permasalahan
    - `GET /permasalahan/{id}` -- View permaslahan by id masalah
    - `GET /permasalahan/{kode_opd}/{tahun}` -- View all permaslahan by kode_opd dan tahun
    - `GET /permasalahan/{kode_opd}/{tahun}/{id_pokin}` -- View all permaslahan by kode_opd, tahun dan id_pokin
    - `POST /permasalahan` -- Save Permasalahan
    - `PATCH /permasalahan/{id}` -- Update Permasalahan
    - `POST /permasalahan/{id}/pilih_permasalahan` -- Pilih permasalahan by id permasalahan dan jadi permasalahan terpilih
    - `DELETE /permasalahan/{id}/hapus_permasalahan_terpilih` -- Delete permasalahan terpilih
      Permasalahan terpilih
    - `GET /permasalahan_terpilih/{kode_opd}/{tahun}` -- View all permaslahan terpilih by kode_opd dan tahun
      Isu strategis
    - `GET /isu_strategis/{id}` -- View isu strategis opd by id isu strategis
    - `GET /isu_strategis/{kode_opd}/{tahun}` -- View all isu strategis opd by id kode_opd dan tahun
    - `POST /isu_strategis` -- Save isu strategis opd
    - `PATCH /isu_strategis/{id}` -- Update isu strategis opd
    - `DELETE /isu_strategis/{id}` -- Delete isu strategis opd
  - [ ] JSON format [lihat request-response-format.md](request-response-format.md)
  - [ ] Error handling and codes:
    - semua field pada request body kecuali id_pokin dan pokin wajib terisi
    - beri keterangan jika field kosong / belum terisi

- [ ] **Identify Service Ownership**

  - [ ] Which microservice handles what
    - isu-permasalahan-service: responsible for gathering isu dan permasalahan opd data
    - perencanaan-service: responsible for gathering data pokin opd
  - [ ] What logic stays within this service:
    - Logic permasalahan dan levelnya
    - Logic pilih permasalahan
    - Logic isu strategis
  - [ ] Delegations to other services (if any):
    - perencanaan-service: handle data pokin

- [ ] **List Integration Points**

  - [ ] Other services involved:
    - perencanaan-service
    - pegawai-service
  - [ ] External APIs or DBs queried
    - none
  - [ ] Auth/permissions service
    - di api gateway

- [ ] **Describe Data Schema Changes**

  - [ ] Table -> cek readme untuk diagram, dbdiagram untuk dbml
  - [ ] Migration strategy:
    - buat database baru
    - gunakan flywaydb untuk java / go
    - gunakan golang migrate untuk go

- [ ] **Add Testing Notes**

  - [ ] Unit test cases
    - Test logic pilih permasalahan
    - Test logic permasalahan terpilih dihapus / batal pilih
  - [ ] Integration test cases
    - none
  - [ ] E2E test cases
    - none
  - [ ] Manual QA checklist:
    - pastikan permasalahan terpilih hanya dapat digunakan sekali di isu strategis opd

- [ ] **Add Status Tracking**
  - [x] Assigned to developer
  - [ ] Reviewed by tech lead
  - [ ] Approved by product owner
