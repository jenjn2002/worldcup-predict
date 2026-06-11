"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { getUser, logout } from "../lib/api";

export default function NavBar() {
  const [user, setUser] = useState(null);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    setUser(getUser());
  }, [pathname]);

  function handleLogout() {
    logout();
    setUser(null);
    router.push("/login");
  }

  return (
    <div className="navbar">
      <Link href="/" className="brand">
        ⚽ World Cup Predictor
      </Link>
      <nav>
        <Link href="/">Trận đấu</Link>
        <Link href="/leaderboard">Bảng xếp hạng</Link>
        {user ? (
          <>
            <span className="tag">👤 {user.username}</span>
            <button onClick={handleLogout}>Đăng xuất</button>
          </>
        ) : (
          <>
            <Link href="/login">Đăng nhập</Link>
            <Link href="/register">Đăng ký</Link>
          </>
        )}
      </nav>
    </div>
  );
}
