"use client";

import { useEffect, useState } from "react";
import { apiFetch, getUser } from "../lib/api";

function formatTime(iso) {
  const d = new Date(iso);
  return d.toLocaleString("vi-VN", {
    weekday: "short",
    day: "2-digit",
    month: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function MatchCard({ match, isLoggedIn, onSaved }) {
  const [home, setHome] = useState(match.my_pred_home ?? "");
  const [away, setAway] = useState(match.my_pred_away ?? "");
  const [saving, setSaving] = useState(false);
  const [msg, setMsg] = useState("");

  const isOpen = match.status === "scheduled" && new Date(match.match_time) > new Date();

  async function handleSave() {
    setMsg("");
    if (home === "" || away === "") {
      setMsg("Vui lòng nhập đủ tỉ số.");
      return;
    }
    setSaving(true);
    try {
      await apiFetch("/api/predictions", {
        method: "POST",
        body: JSON.stringify({
          match_id: match.id,
          pred_home: Number(home),
          pred_away: Number(away),
        }),
      });
      setMsg("Đã lưu dự đoán!");
      onSaved?.();
    } catch (err) {
      setMsg(err.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="card match-card">
      <div className="match-header">
        <span className="tag">
          {match.stage === "group" ? `Bảng ${match.group_name}` : match.stage}
        </span>
        <span>{formatTime(match.match_time)}</span>
      </div>

      <div className="match-teams">
        <span className="team-name">{match.home_team}</span>
        {match.status === "finished" ? (
          <span className="final-score">
            {match.home_score} - {match.away_score}
          </span>
        ) : (
          <span className="muted">vs</span>
        )}
        <span className="team-name">{match.away_team}</span>
      </div>

      {isLoggedIn ? (
        <div className="predict-row">
          <div className="score-input">
            <input
              type="number"
              min={0}
              value={home}
              disabled={!isOpen}
              onChange={(e) => setHome(e.target.value)}
            />
            <span>:</span>
            <input
              type="number"
              min={0}
              value={away}
              disabled={!isOpen}
              onChange={(e) => setAway(e.target.value)}
            />
          </div>
          {isOpen ? (
            <button className="btn" onClick={handleSave} disabled={saving}>
              {saving ? "Đang lưu..." : "Lưu dự đoán"}
            </button>
          ) : (
            match.my_points !== undefined &&
            match.my_points !== null &&
            match.status === "finished" && (
              <span className={`points-tag points-${match.my_points}`}>
                +{match.my_points} điểm
              </span>
            )
          )}
          {!isOpen && match.status !== "finished" && (
            <span className="muted">Đã đóng dự đoán</span>
          )}
        </div>
      ) : (
        <p className="muted" style={{ textAlign: "center" }}>
          Đăng nhập để dự đoán trận này
        </p>
      )}

      {msg && <p className="muted">{msg}</p>}
    </div>
  );
}

export default function HomePage() {
  const [matches, setMatches] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [user, setUserState] = useState(null);

  function load() {
    setLoading(true);
    apiFetch("/api/matches")
      .then(setMatches)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    setUserState(getUser());
    load();
  }, []);

  return (
    <div>
      <h1>⚽ Lịch thi đấu &amp; Dự đoán</h1>
      <p className="muted">
        Đoán đúng tỉ số chính xác = 3 điểm. Đoán đúng kết quả (thắng/thua/hòa) = 1 điểm.
      </p>
      {error && <div className="error">{error}</div>}
      {loading ? (
        <p className="muted">Đang tải...</p>
      ) : (
        matches.map((m) => (
          <MatchCard key={m.id} match={m} isLoggedIn={!!user} onSaved={load} />
        ))
      )}
      {!loading && matches.length === 0 && <p className="muted">Chưa có trận đấu nào.</p>}
    </div>
  );
}
