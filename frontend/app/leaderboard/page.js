"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "../../lib/api";

export default function LeaderboardPage() {
  const [entries, setEntries] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiFetch("/api/leaderboard")
      .then(setEntries)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <h1>🏆 Bảng xếp hạng</h1>
      {error && <div className="error">{error}</div>}
      {loading ? (
        <p className="muted">Đang tải...</p>
      ) : (
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>#</th>
                <th>Người chơi</th>
                <th>Tổng điểm</th>
                <th>Đoán đúng tỉ số</th>
                <th>Số lượt đoán</th>
              </tr>
            </thead>
            <tbody>
              {entries.map((e, i) => (
                <tr key={e.user_id}>
                  <td className={i === 0 ? "rank-1" : ""}>{i + 1}</td>
                  <td className={i === 0 ? "rank-1" : ""}>{e.username}</td>
                  <td>{e.total_points}</td>
                  <td>{e.exact_count}</td>
                  <td>{e.predicted_count}</td>
                </tr>
              ))}
              {entries.length === 0 && (
                <tr>
                  <td colSpan={5} className="muted">
                    Chưa có dữ liệu.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
      <p className="muted">
        Tính điểm: đoán đúng tỉ số chính xác = 3 điểm, đoán đúng kết quả thắng/thua/hòa = 1 điểm.
      </p>
    </div>
  );
}
