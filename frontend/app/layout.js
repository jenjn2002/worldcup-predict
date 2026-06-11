import "./globals.css";
import NavBar from "./NavBar";

export const metadata = {
  title: "World Cup Predictor",
  description: "Đoán tỉ số các trận đấu World Cup",
};

export default function RootLayout({ children }) {
  return (
    <html lang="vi">
      <body>
        <NavBar />
        <div className="container">{children}</div>
      </body>
    </html>
  );
}
