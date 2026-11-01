CREATE TABLE IF NOT EXISTS authorized_people (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	student_id INTEGER NOT NULL REFERENCES students(id),
	full_name TEXT NOT NULL,
	document_type TEXT NOT NULL CHECK (document_type IN ('CC', 'TI', 'CE', 'PASSPORT', 'OTHER')),
	document_number TEXT NOT NULL,
	relationship TEXT NOT NULL,
	phone TEXT NOT NULL,
	UNIQUE(student_id, document_type, document_number)
);
CREATE INDEX IF NOT EXISTS idx_authorized_people_student_id ON authorized_people(student_id);
CREATE INDEX IF NOT EXISTS idx_students_full_name ON students(full_name);
