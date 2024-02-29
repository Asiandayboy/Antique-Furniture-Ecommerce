import { useEffect, useState } from "react"


const MAX_PRICE_RANGE = 1_000_000

type PriceRange = {
  priceMin: number,
  priceMax: number,
}


type SearchQuery = {
  searchString: string,
  type: string,
  material: string,
  condition: string,
  style: string,
  priceRange: PriceRange,
}



export default function MarketFilter() {
  const [priceInput, setPriceInput] = useState<PriceRange>({
    priceMin: 0,
    priceMax: MAX_PRICE_RANGE/2
  })
  const [searchInput, setSearchInput] = useState<SearchQuery>({
    searchString: "",
    type: "All",
    material: "All",
    condition: "All",
    style: "All",
    priceRange: priceInput
  })




  useEffect(() => {
    console.log(searchInput)
  }, [searchInput])


  return (
    <div className="market-filter_wrapper">
      <div className="filters_wrapper">
        <div className="select_filters">
          <div className="furniture-type_filter">
            <label htmlFor="">Furniture Type </label>
            <select onChange={
              (e) => setSearchInput({...searchInput, type: e.currentTarget.value})
            } name="" id="">
              <option value="All">All</option>
              <option value="Bed">Bed</option>
              <option value="Table">Table</option>
              <option value="Chair">Chair</option>
              <option value="Nightstand">Nightstand</option>
              <option value="Desk">Desk</option>
              <option value="Lamp">Lamp</option>
            </select>
          </div>

          <div className="material_filter">
            <label htmlFor="">Material </label>
            <select onChange={
              e => setSearchInput({...searchInput, material: e.currentTarget.value})
            } name="" id="">
              <option value="All">All</option>
              <option value="TigerMaple">Tiger Maple</option>
              <option value="Walnut">Walnut</option>
              <option value="Oak">Oak</option>
              <option value="Cherry">Cherry</option>
              <option value="Mahagony">Mahagony</option>
              <option value="Maple">Maple</option>
            </select>
          </div>

          <div className="condition_filter">
            <label htmlFor="">Condition </label>
            <select onChange={
              e => setSearchInput({...searchInput, condition: e.currentTarget.value})
            }  name="" id="">
              <option value="All">All</option>
              <option value="Mint">Mint</option>
              <option value="Excellent">Excellent</option>
              <option value="Good">Good</option>
              <option value="Worn">Worn</option>
              <option value="Restored">Restored</option>
              <option value="Original Finish">Original Finish</option>
            </select>
          </div>

          <div className="style_filter">
            <label htmlFor="">Style </label>
            <select onChange={
              e => setSearchInput({...searchInput, style: e.currentTarget.value})
            }  name="" id="">
              <option value="All">All</option>
              <option value="Victorian">Victorian</option>
              <option value="English">English</option>
              <option value="Baroque">Baroque</option>
              <option value="Federal">Federal</option>
              <option value="Rococo">Rococo</option>
              <option value="Sheraton">Sheraton</option>
            </select>
          </div>
        </div>

        <div className="price-range_filter">
          <div>
            <label htmlFor="">Price Range </label>
            <div className="price-ranges">
              <div>
                <label htmlFor="">min</label>
                <input onChange={
                  (e) => {
                    setPriceInput({...priceInput, priceMin: +e.currentTarget.value})
                    setSearchInput({...searchInput, priceRange: priceInput})
                  }
                } type="range" min="0" max="1000000" value={priceInput.priceMin}/>
                <div className="price_text">${priceInput.priceMin}</div>
              </div>
              <div>
                <label htmlFor="">max</label>
                <input onChange={
                  (e) => {
                    setPriceInput({...priceInput, priceMax: +e.currentTarget.value})
                    setSearchInput({...searchInput, priceRange: priceInput})
                  }
                } type="range" min="0" max="1000000" value={priceInput.priceMax}/>
                <div className="price_text">${priceInput.priceMax}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="search-bar">
        <label htmlFor="">Search for furniture:</label>
        <div>
          <input onChange={
            (e) => setSearchInput({...searchInput, searchString: e.currentTarget.value})
          } type="text" />
          <button>Search</button>
        </div>
      </div>

    </div>
  )
}