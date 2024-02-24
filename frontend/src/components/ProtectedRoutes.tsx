import { Navigate, Outlet } from "react-router-dom"


type Props = {
  auth: boolean
  redirect: string
}


export default function ProtectedRoutes({ auth, redirect }: Props) {
  return (
    auth ? <Outlet /> : <Navigate to={redirect}/>
  )
}
