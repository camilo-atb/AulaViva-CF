CREATE TABLE IF NOT EXISTS courses (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	course_id INTEGER NOT NULL REFERENCES courses(id),
	full_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	course_id INTEGER NOT NULL REFERENCES courses(id),
	name TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	event_date TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS authorizations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id INTEGER NOT NULL REFERENCES events(id),
	student_id INTEGER NOT NULL REFERENCES students(id),
	request_token TEXT NOT NULL UNIQUE,
	status TEXT NOT NULL CHECK (status IN ('pending', 'authorized', 'not_authorized')),
	responded_at TEXT,
	UNIQUE(event_id, student_id),
	CHECK ((status = 'pending' AND responded_at IS NULL) OR (status <> 'pending' AND responded_at IS NOT NULL))
);
CREATE INDEX IF NOT EXISTS idx_authorizations_event_id ON authorizations(event_id);
CREATE INDEX IF NOT EXISTS idx_authorizations_student_id ON authorizations(student_id);
CREATE INDEX IF NOT EXISTS idx_authorizations_status ON authorizations(status);
