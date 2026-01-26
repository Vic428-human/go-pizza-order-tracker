import { useEffect, useState } from "react";

export default function AdminDashboard() {
  const [data, setData] = useState(null);

  useEffect(() => {
    // gin-contrib/sessions，React 呼叫 API 時要帶 cookie。
    fetch("/api/admin/dashboard", { credentials: "include" })
      .then(res => res.json())
      .then(setData);
  }, []);

  if (!data) return <p>Loading...</p>;

  return (
    <div>
      <p>Welcome, {data.username}</p>
      <p>Status: {data.status}</p>
      <ul>
        {data.orders.map((o, i) => <li key={i}>{o}</li>)}
      </ul>
    </div>
  );
}
