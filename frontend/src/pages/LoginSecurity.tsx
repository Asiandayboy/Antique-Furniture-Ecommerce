import Navbar from "../components/Navbar";
import { useAccountDataContext } from "../contexts/accountDataContext";

export default function LoginSecurity() {

  const { userData } = useAccountDataContext()

  return (
    <>
      <Navbar />
      <div>
        <h1>LoginSecurity</h1>
        <div>
          <div>Email: {userData?.email}</div>
          <div>Phone: {userData?.phone}</div>
          <div>Password: ----------</div>
        </div>

      </div>
    </>
  )
}