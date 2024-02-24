import Navbar from "../components/Navbar"

type Props = {
  isLoggedIn: boolean
}


export default function Home({ isLoggedIn }: Props) {
  return (
    <>
      <Navbar isLoggedIn={isLoggedIn} />
      <main>
        <h1>Shop for great antique furniture!</h1>
      </main>
    </>
  )
}