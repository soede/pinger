package repository

const (
	//getAll = `SELECT * FROM containers LIMIT $1 OFFSET $2`
	getAll = `WITH RankedContainers AS (
            SELECT
                container_id,
                container_status,
                addr,
                p_duration,
                pinged_at,
                ROW_NUMBER() OVER (PARTITION BY addr ORDER BY pinged_at DESC) AS row_num
            FROM containers
        )
        SELECT
            container_id,
            container_status,
            addr,
            p_duration,
            pinged_at
        FROM RankedContainers
        WHERE row_num = 1
        LIMIT $1 OFFSET $2;`

	getHistory = `SELECT * FROM containers LIMIT $1 OFFSET $2`

	countQuery = `SELECT COUNT(DISTINCT addr) FROM containers`

	countHistoryQuery = `SELECT COUNT(*) FROM containers`

	setAll = `INSERT INTO containers (container_status, addr, pinged_at, p_duration) 
              VALUES (:container_status, :addr, :pinged_at, :p_duration)`

	getByIP     = `SELECT * FROM containers WHERE addr LIKE $1`
	getByStatus = `SELECT * FROM containers WHERE container_status LIKE $1`
)
