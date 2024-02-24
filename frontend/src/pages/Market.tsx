import Navbar from "../components/Navbar"

type Props = {
  isLoggedIn: boolean
}

export default function Market({ isLoggedIn }: Props) {
  return (
    <>
      <Navbar />
      <div>Market</div>
    </>
  )
}