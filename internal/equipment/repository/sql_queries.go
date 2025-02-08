package repository

const (
	qCreateEquipment = `INSERT INTO equipment (name, short_description, full_description) VALUES ($1, $2, $3) RETURNING *`
	qUpdateEquipment = `UPDATE equipment SET name=$1, short_description=$2, full_description=$3 WHERE equipment_id=$4`
	qDeleteEquipment = `DELETE FROM equipment WHERE equipment_id = $1`
	qGetEquipment    = `SELECT equipment_id, name, short_description, full_description FROM equipment WHERE equipment_id = $1`

	qGetTotal               = `SELECT COUNT(equipment_id) FROM equipment`
	qGetTotalReservedByUser = `SELECT COUNT(equipment_id) 
								FROM (
									SELECT DISTINCT equipment_id
									FROM usersEquipment
									WHERE user_id = $1 AND CURRENT_TIMESTAMP < reservation_end
								)`

	qGetEquipments = `WITH nearest_reservations AS (
		SELECT DISTINCT  equipment_id, MIN(reservation_start) AS reservation_start, MIN(reservation_end) AS reservation_end
		FROM usersEquipment 
		WHERE CURRENT_TIMESTAMP < reservation_end
		GROUP BY equipment_id
	)
	SELECT equipment_id, name, short_description,
		CASE WHEN reservation_start IS NULL THEN false ELSE CURRENT_TIMESTAMP BETWEEN reservation_start AND reservation_end END AS reserved
	FROM equipment
	LEFT JOIN nearest_reservations using(equipment_id)
	ORDER BY reserved
	OFFSET $1 
	LIMIT $2`

	// qGetEquipments = `SELECT equipment_id, name, short_description
	// 				 FROM equipment
	// 				 ORDER BY COALESCE(NULLIF($1, ''), name) OFFSET $2 LIMIT $3`
	qGetUserEquipments = `SELECT DISTINCT equipment_id, name, short_description, true AS reserved
	FROM equipment 
	INNER JOIN usersEquipment using(equipment_id)
	WHERE user_id = $3 AND CURRENT_TIMESTAMP < reservation_end
	OFFSET $1 LIMIT $2`

	// qGetReservationInfo = `SELECT reservation_start, reservation_end
	// 					FROM usersEquipment
	// 					WHERE equipment_id = $1
	// 					AND CURRENT_TIMESTAMP BETWEEN reservation_start AND reservation_end`

	qGetReservationInfo = `SELECT reservation_start, reservation_end
	FROM usersEquipment
	WHERE equipment_id = $1
	ORDER BY reservation_start`

	qIsReserved = `SELECT true AS reserved
					FROM usersEquipment
					WHERE equipment_id = $1
					AND ($2 BETWEEN reservation_start AND reservation_end
					OR $3 BETWEEN reservation_start AND reservation_end)
					LIMIT 1`

	qReserve = `INSERT INTO usersEquipment (user_id, equipment_id, reservation_start, reservation_end)
				VALUES ($1, $2, $3, $4)`
)
