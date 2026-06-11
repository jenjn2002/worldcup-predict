-- =========================================================
-- World Cup Prediction Game - Database Schema
-- =========================================================

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    is_admin      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS teams (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    group_name VARCHAR(5)
);

CREATE TABLE IF NOT EXISTS matches (
    id            SERIAL PRIMARY KEY,
    stage         VARCHAR(30) NOT NULL DEFAULT 'group', -- group, round_of_32, round_of_16, quarter, semi, third_place, final
    group_name    VARCHAR(5),
    home_team_id  INTEGER REFERENCES teams(id),
    away_team_id  INTEGER REFERENCES teams(id),
    match_time    TIMESTAMPTZ NOT NULL,
    home_score    INTEGER,
    away_score    INTEGER,
    status        VARCHAR(20) NOT NULL DEFAULT 'scheduled' -- scheduled, finished
);

CREATE TABLE IF NOT EXISTS predictions (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id    INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    pred_home   INTEGER NOT NULL,
    pred_away   INTEGER NOT NULL,
    points      INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, match_id)
);

CREATE INDEX IF NOT EXISTS idx_matches_time ON matches(match_time);
CREATE INDEX IF NOT EXISTS idx_predictions_user ON predictions(user_id);
CREATE INDEX IF NOT EXISTS idx_predictions_match ON predictions(match_id);

-- =========================================================
-- Seed data: 12 groups x 4 teams (World Cup 2026 group stage)
-- Note: groups C, E, F, H, J, K, L use placeholder team names
-- (TBD ..) - update them with real team names once you know
-- the full draw, by editing the `teams` table.
-- =========================================================
-- TEAMS
INSERT INTO teams (name, group_name) VALUES
('Mexico', 'A'),
('South Africa', 'A'),
('South Korea', 'A'),
('Czechia', 'A'),
('Canada', 'B'),
('Bosnia and Herzegovina', 'B'),
('Qatar', 'B'),
('Switzerland', 'B'),
('TBD C1', 'C'),
('TBD C2', 'C'),
('TBD C3', 'C'),
('TBD C4', 'C'),
('United States', 'D'),
('Paraguay', 'D'),
('Australia', 'D'),
('Turkiye', 'D'),
('TBD E1', 'E'),
('TBD E2', 'E'),
('TBD E3', 'E'),
('TBD E4', 'E'),
('TBD F1', 'F'),
('TBD F2', 'F'),
('TBD F3', 'F'),
('TBD F4', 'F'),
('Belgium', 'G'),
('Egypt', 'G'),
('Iran', 'G'),
('New Zealand', 'G'),
('TBD H1', 'H'),
('TBD H2', 'H'),
('TBD H3', 'H'),
('TBD H4', 'H'),
('France', 'I'),
('Senegal', 'I'),
('Iraq', 'I'),
('Norway', 'I'),
('TBD J1', 'J'),
('TBD J2', 'J'),
('TBD J3', 'J'),
('TBD J4', 'J'),
('TBD K1', 'K'),
('TBD K2', 'K'),
('TBD K3', 'K'),
('TBD K4', 'K'),
('TBD L1', 'L'),
('TBD L2', 'L'),
('TBD L3', 'L'),
('TBD L4', 'L');

-- MATCHES (group stage, round robin within each group)
INSERT INTO matches (stage, group_name, home_team_id, away_team_id, match_time) VALUES
('group', 'A', (SELECT id FROM teams WHERE name='Mexico'), (SELECT id FROM teams WHERE name='South Africa'), '2026-06-11 15:00:00+00'),
('group', 'A', (SELECT id FROM teams WHERE name='Mexico'), (SELECT id FROM teams WHERE name='South Korea'), '2026-06-11 18:00:00+00'),
('group', 'A', (SELECT id FROM teams WHERE name='Mexico'), (SELECT id FROM teams WHERE name='Czechia'), '2026-06-11 21:00:00+00'),
('group', 'A', (SELECT id FROM teams WHERE name='South Africa'), (SELECT id FROM teams WHERE name='South Korea'), '2026-06-12 00:00:00+00'),
('group', 'A', (SELECT id FROM teams WHERE name='South Africa'), (SELECT id FROM teams WHERE name='Czechia'), '2026-06-12 03:00:00+00'),
('group', 'A', (SELECT id FROM teams WHERE name='South Korea'), (SELECT id FROM teams WHERE name='Czechia'), '2026-06-12 06:00:00+00'),
('group', 'B', (SELECT id FROM teams WHERE name='Canada'), (SELECT id FROM teams WHERE name='Bosnia and Herzegovina'), '2026-06-12 09:00:00+00'),
('group', 'B', (SELECT id FROM teams WHERE name='Canada'), (SELECT id FROM teams WHERE name='Qatar'), '2026-06-12 12:00:00+00'),
('group', 'B', (SELECT id FROM teams WHERE name='Canada'), (SELECT id FROM teams WHERE name='Switzerland'), '2026-06-12 15:00:00+00'),
('group', 'B', (SELECT id FROM teams WHERE name='Bosnia and Herzegovina'), (SELECT id FROM teams WHERE name='Qatar'), '2026-06-12 18:00:00+00'),
('group', 'B', (SELECT id FROM teams WHERE name='Bosnia and Herzegovina'), (SELECT id FROM teams WHERE name='Switzerland'), '2026-06-12 21:00:00+00'),
('group', 'B', (SELECT id FROM teams WHERE name='Qatar'), (SELECT id FROM teams WHERE name='Switzerland'), '2026-06-13 00:00:00+00'),
('group', 'C', (SELECT id FROM teams WHERE name='TBD C1'), (SELECT id FROM teams WHERE name='TBD C2'), '2026-06-13 03:00:00+00'),
('group', 'C', (SELECT id FROM teams WHERE name='TBD C1'), (SELECT id FROM teams WHERE name='TBD C3'), '2026-06-13 06:00:00+00'),
('group', 'C', (SELECT id FROM teams WHERE name='TBD C1'), (SELECT id FROM teams WHERE name='TBD C4'), '2026-06-13 09:00:00+00'),
('group', 'C', (SELECT id FROM teams WHERE name='TBD C2'), (SELECT id FROM teams WHERE name='TBD C3'), '2026-06-13 12:00:00+00'),
('group', 'C', (SELECT id FROM teams WHERE name='TBD C2'), (SELECT id FROM teams WHERE name='TBD C4'), '2026-06-13 15:00:00+00'),
('group', 'C', (SELECT id FROM teams WHERE name='TBD C3'), (SELECT id FROM teams WHERE name='TBD C4'), '2026-06-13 18:00:00+00'),
('group', 'D', (SELECT id FROM teams WHERE name='United States'), (SELECT id FROM teams WHERE name='Paraguay'), '2026-06-13 21:00:00+00'),
('group', 'D', (SELECT id FROM teams WHERE name='United States'), (SELECT id FROM teams WHERE name='Australia'), '2026-06-14 00:00:00+00'),
('group', 'D', (SELECT id FROM teams WHERE name='United States'), (SELECT id FROM teams WHERE name='Turkiye'), '2026-06-14 03:00:00+00'),
('group', 'D', (SELECT id FROM teams WHERE name='Paraguay'), (SELECT id FROM teams WHERE name='Australia'), '2026-06-14 06:00:00+00'),
('group', 'D', (SELECT id FROM teams WHERE name='Paraguay'), (SELECT id FROM teams WHERE name='Turkiye'), '2026-06-14 09:00:00+00'),
('group', 'D', (SELECT id FROM teams WHERE name='Australia'), (SELECT id FROM teams WHERE name='Turkiye'), '2026-06-14 12:00:00+00'),
('group', 'E', (SELECT id FROM teams WHERE name='TBD E1'), (SELECT id FROM teams WHERE name='TBD E2'), '2026-06-14 15:00:00+00'),
('group', 'E', (SELECT id FROM teams WHERE name='TBD E1'), (SELECT id FROM teams WHERE name='TBD E3'), '2026-06-14 18:00:00+00'),
('group', 'E', (SELECT id FROM teams WHERE name='TBD E1'), (SELECT id FROM teams WHERE name='TBD E4'), '2026-06-14 21:00:00+00'),
('group', 'E', (SELECT id FROM teams WHERE name='TBD E2'), (SELECT id FROM teams WHERE name='TBD E3'), '2026-06-15 00:00:00+00'),
('group', 'E', (SELECT id FROM teams WHERE name='TBD E2'), (SELECT id FROM teams WHERE name='TBD E4'), '2026-06-15 03:00:00+00'),
('group', 'E', (SELECT id FROM teams WHERE name='TBD E3'), (SELECT id FROM teams WHERE name='TBD E4'), '2026-06-15 06:00:00+00'),
('group', 'F', (SELECT id FROM teams WHERE name='TBD F1'), (SELECT id FROM teams WHERE name='TBD F2'), '2026-06-15 09:00:00+00'),
('group', 'F', (SELECT id FROM teams WHERE name='TBD F1'), (SELECT id FROM teams WHERE name='TBD F3'), '2026-06-15 12:00:00+00'),
('group', 'F', (SELECT id FROM teams WHERE name='TBD F1'), (SELECT id FROM teams WHERE name='TBD F4'), '2026-06-15 15:00:00+00'),
('group', 'F', (SELECT id FROM teams WHERE name='TBD F2'), (SELECT id FROM teams WHERE name='TBD F3'), '2026-06-15 18:00:00+00'),
('group', 'F', (SELECT id FROM teams WHERE name='TBD F2'), (SELECT id FROM teams WHERE name='TBD F4'), '2026-06-15 21:00:00+00'),
('group', 'F', (SELECT id FROM teams WHERE name='TBD F3'), (SELECT id FROM teams WHERE name='TBD F4'), '2026-06-16 00:00:00+00'),
('group', 'G', (SELECT id FROM teams WHERE name='Belgium'), (SELECT id FROM teams WHERE name='Egypt'), '2026-06-16 03:00:00+00'),
('group', 'G', (SELECT id FROM teams WHERE name='Belgium'), (SELECT id FROM teams WHERE name='Iran'), '2026-06-16 06:00:00+00'),
('group', 'G', (SELECT id FROM teams WHERE name='Belgium'), (SELECT id FROM teams WHERE name='New Zealand'), '2026-06-16 09:00:00+00'),
('group', 'G', (SELECT id FROM teams WHERE name='Egypt'), (SELECT id FROM teams WHERE name='Iran'), '2026-06-16 12:00:00+00'),
('group', 'G', (SELECT id FROM teams WHERE name='Egypt'), (SELECT id FROM teams WHERE name='New Zealand'), '2026-06-16 15:00:00+00'),
('group', 'G', (SELECT id FROM teams WHERE name='Iran'), (SELECT id FROM teams WHERE name='New Zealand'), '2026-06-16 18:00:00+00'),
('group', 'H', (SELECT id FROM teams WHERE name='TBD H1'), (SELECT id FROM teams WHERE name='TBD H2'), '2026-06-16 21:00:00+00'),
('group', 'H', (SELECT id FROM teams WHERE name='TBD H1'), (SELECT id FROM teams WHERE name='TBD H3'), '2026-06-17 00:00:00+00'),
('group', 'H', (SELECT id FROM teams WHERE name='TBD H1'), (SELECT id FROM teams WHERE name='TBD H4'), '2026-06-17 03:00:00+00'),
('group', 'H', (SELECT id FROM teams WHERE name='TBD H2'), (SELECT id FROM teams WHERE name='TBD H3'), '2026-06-17 06:00:00+00'),
('group', 'H', (SELECT id FROM teams WHERE name='TBD H2'), (SELECT id FROM teams WHERE name='TBD H4'), '2026-06-17 09:00:00+00'),
('group', 'H', (SELECT id FROM teams WHERE name='TBD H3'), (SELECT id FROM teams WHERE name='TBD H4'), '2026-06-17 12:00:00+00'),
('group', 'I', (SELECT id FROM teams WHERE name='France'), (SELECT id FROM teams WHERE name='Senegal'), '2026-06-17 15:00:00+00'),
('group', 'I', (SELECT id FROM teams WHERE name='France'), (SELECT id FROM teams WHERE name='Iraq'), '2026-06-17 18:00:00+00'),
('group', 'I', (SELECT id FROM teams WHERE name='France'), (SELECT id FROM teams WHERE name='Norway'), '2026-06-17 21:00:00+00'),
('group', 'I', (SELECT id FROM teams WHERE name='Senegal'), (SELECT id FROM teams WHERE name='Iraq'), '2026-06-18 00:00:00+00'),
('group', 'I', (SELECT id FROM teams WHERE name='Senegal'), (SELECT id FROM teams WHERE name='Norway'), '2026-06-18 03:00:00+00'),
('group', 'I', (SELECT id FROM teams WHERE name='Iraq'), (SELECT id FROM teams WHERE name='Norway'), '2026-06-18 06:00:00+00'),
('group', 'J', (SELECT id FROM teams WHERE name='TBD J1'), (SELECT id FROM teams WHERE name='TBD J2'), '2026-06-18 09:00:00+00'),
('group', 'J', (SELECT id FROM teams WHERE name='TBD J1'), (SELECT id FROM teams WHERE name='TBD J3'), '2026-06-18 12:00:00+00'),
('group', 'J', (SELECT id FROM teams WHERE name='TBD J1'), (SELECT id FROM teams WHERE name='TBD J4'), '2026-06-18 15:00:00+00'),
('group', 'J', (SELECT id FROM teams WHERE name='TBD J2'), (SELECT id FROM teams WHERE name='TBD J3'), '2026-06-18 18:00:00+00'),
('group', 'J', (SELECT id FROM teams WHERE name='TBD J2'), (SELECT id FROM teams WHERE name='TBD J4'), '2026-06-18 21:00:00+00'),
('group', 'J', (SELECT id FROM teams WHERE name='TBD J3'), (SELECT id FROM teams WHERE name='TBD J4'), '2026-06-19 00:00:00+00'),
('group', 'K', (SELECT id FROM teams WHERE name='TBD K1'), (SELECT id FROM teams WHERE name='TBD K2'), '2026-06-19 03:00:00+00'),
('group', 'K', (SELECT id FROM teams WHERE name='TBD K1'), (SELECT id FROM teams WHERE name='TBD K3'), '2026-06-19 06:00:00+00'),
('group', 'K', (SELECT id FROM teams WHERE name='TBD K1'), (SELECT id FROM teams WHERE name='TBD K4'), '2026-06-19 09:00:00+00'),
('group', 'K', (SELECT id FROM teams WHERE name='TBD K2'), (SELECT id FROM teams WHERE name='TBD K3'), '2026-06-19 12:00:00+00'),
('group', 'K', (SELECT id FROM teams WHERE name='TBD K2'), (SELECT id FROM teams WHERE name='TBD K4'), '2026-06-19 15:00:00+00'),
('group', 'K', (SELECT id FROM teams WHERE name='TBD K3'), (SELECT id FROM teams WHERE name='TBD K4'), '2026-06-19 18:00:00+00'),
('group', 'L', (SELECT id FROM teams WHERE name='TBD L1'), (SELECT id FROM teams WHERE name='TBD L2'), '2026-06-19 21:00:00+00'),
('group', 'L', (SELECT id FROM teams WHERE name='TBD L1'), (SELECT id FROM teams WHERE name='TBD L3'), '2026-06-20 00:00:00+00'),
('group', 'L', (SELECT id FROM teams WHERE name='TBD L1'), (SELECT id FROM teams WHERE name='TBD L4'), '2026-06-20 03:00:00+00'),
('group', 'L', (SELECT id FROM teams WHERE name='TBD L2'), (SELECT id FROM teams WHERE name='TBD L3'), '2026-06-20 06:00:00+00'),
('group', 'L', (SELECT id FROM teams WHERE name='TBD L2'), (SELECT id FROM teams WHERE name='TBD L4'), '2026-06-20 09:00:00+00'),
('group', 'L', (SELECT id FROM teams WHERE name='TBD L3'), (SELECT id FROM teams WHERE name='TBD L4'), '2026-06-20 12:00:00+00');
