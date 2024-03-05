import { Navigate, Outlet } from "react-router-dom"


type Props = {
  auth: boolean
  redirect: string
}

/**This component renders child components if auth is true, else
 * it navigates to the component provided as redirect
 */
export default function ProtectedRoutes({ auth, redirect }: Props) {
  return (
    auth ? <Outlet /> : <Navigate to={redirect}/>
  )
}
