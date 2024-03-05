import { useEffect, useState, useRef } from "react"
import { FurnitureListing } from "../pages/Market"
import * as FurnitureTypes from "../types/furnitureTypes"
import { deepEqual } from "../util/obj"


const MAX_PRICE_RANGE = 100_000


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

type SearchCategories = {
  type: FurnitureTypes.FurnitureType[],
  style: FurnitureTypes.FurnitureStyle[],
  material: FurnitureTypes.FurnitureMaterial[]
  bedSize: FurnitureTypes.FurnitureBedSize[]
}


type Props = {
  dataSet: FurnitureListing[]
  setDataSet: React.Dispatch<React.SetStateAction<FurnitureListing[]>>
}



function createSearchCategories(words: string[]): SearchCategories {
  const categories: SearchCategories = {
    type: [],
    style: [],
    material: [],
    bedSize: []
  }

  words.forEach((word) => {
    if (FurnitureTypes.isFurnitureBedSize(word)) {
      categories.bedSize.push(word)
    } else if (FurnitureTypes.isFurnitureMaterial(word)) {
      categories.material.push(word)
    } else if (FurnitureTypes.isFurnitureStyle(word)) {
      categories.style.push(word)
    } else if (FurnitureTypes.isFurnitureType(word)) {
      categories.type.push(word)
    }
  })

  return categories
}

const defaultPriceRange: PriceRange = {
  priceMin: 0,
  priceMax: MAX_PRICE_RANGE/2
}

const defaultSearchQuery: SearchQuery = {
  searchString: "",
  type: "All",
  material: "All",
  condition: "All",
  style: "All",
  priceRange: defaultPriceRange
}



export default function MarketFilter({ dataSet, setDataSet }: Props) {
  const [priceInput, setPriceInput] = useState<PriceRange>(defaultPriceRange)
  const [searchQuery, setSearchQuery] = useState<SearchQuery>(defaultSearchQuery)

  const typeRef = useRef<HTMLSelectElement>(null)
  const materialRef = useRef<HTMLSelectElement>(null)
  const conditionRef = useRef<HTMLSelectElement>(null)
  const styleRef = useRef<HTMLSelectElement>(null)
  const priceMinRef = useRef(null)
  const priceMaxRef = useRef(null)


  function clearFilter(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    if (typeRef.current) {
      typeRef.current.selectedIndex = 0
    }
    if (materialRef.current) {
      materialRef.current.selectedIndex = 0
    }
    if (conditionRef.current) {
      conditionRef.current.selectedIndex = 0
    }
    if (styleRef.current) {
      styleRef.current.selectedIndex = 0
    }

    setPriceInput(defaultPriceRange)
    setSearchQuery(defaultSearchQuery)
  }


  function onSearch(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    // search all if filters are set to default, including search bar
    if (deepEqual(searchQuery, defaultSearchQuery)) {
      setDataSet(dataSet)
      return
    }

    const words = searchQuery.searchString
    .toLowerCase()
    .split(" ")
    .filter(word => word !== "")
    
    const categories = createSearchCategories(words)

    console.log(categories)

    setDataSet(dataSet.filter((listing) => {
      let flag = true;

      if (categories.type.length > 0) {
        flag = categories.type.some(word => listing.type.toLowerCase().includes(word))
      } 

      if (categories.style.length > 0) {
        flag = categories.style.some(word => listing.style.toLowerCase().includes(word))
      }

      if (categories.material.length > 0) {
        flag = categories.material.some(word => listing.material.toLowerCase().includes(word))
      }

      if (listing.type.toLowerCase() == "bed" && categories.bedSize.length > 0) {
        flag = categories.bedSize.some(word => listing.title.toLowerCase().includes(word))
      }
      

      // match listing attribute with appropriate filter type
      if (searchQuery.type != "All" 
      && listing.type.toLowerCase() != searchQuery.type) {
        return false
      }

      if (searchQuery.material != "All" 
      && listing.material.toLowerCase() != searchQuery.material) {
        return false
      }

      if (searchQuery.condition != "All" 
      && listing.condition.toLowerCase() != searchQuery.condition) {
        return false
      }

      if (searchQuery.style != "All" 
      && listing.style.toLowerCase() != searchQuery.style) {
        return false
      }

      if (+listing.cost < priceInput.priceMin || +listing.cost > priceInput.priceMax) {
        return false
      }

      return flag

    }))
  }



  return (
    <div className="market-filter_wrapper">
      <div className="filters_wrapper">
        <div className="select_filters">
          <div className="furniture-type_filter">
            <label htmlFor="">Furniture Type </label>
            <select ref={typeRef} onChange={
              (e) => setSearchQuery({...searchQuery, type: e.currentTarget.value})
            } name="" id="">
              <option value="All">All</option>
              <option value="bed">Bed</option>
              <option value="table">Table</option>
              <option value="chair">Chair</option>
              <option value="nightstand">Nightstand</option>
              <option value="desk">Desk</option>
              <option value="lamp">Lamp</option>
            </select>
          </div>

          <div className="material_filter">
            <label htmlFor="">Material </label>
            <select ref={materialRef} onChange={
              e => setSearchQuery({...searchQuery, material: e.currentTarget.value})
            } name="" id="">
              <option value="All">All</option>
              <option value="tiger maple">Tiger Maple</option>
              <option value="walnut">Walnut</option>
              <option value="oak">Oak</option>
              <option value="cherry">Cherry</option>
              <option value="mahagony">Mahagony</option>
              <option value="maple">Maple</option>
              <option value="rosewood">Rosewood</option>
              <option value="birch">Birch</option>
            </select>
          </div>

          <div className="condition_filter">
            <label htmlFor="">Condition </label>
            <select ref={conditionRef} onChange={
              e => setSearchQuery({...searchQuery, condition: e.currentTarget.value})
            }  name="" id="">
              <option value="All">All</option>
              <option value="mint">Mint</option>
              <option value="excellent">Excellent</option>
              <option value="good">Good</option>
              <option value="worn">Worn</option>
              <option value="restored">Restored</option>
              <option value="original Finish">Original Finish</option>
            </select>
          </div>

          <div className="style_filter">
            <label htmlFor="">Style </label>
            <select ref={styleRef} onChange={
              e => setSearchQuery({...searchQuery, style: e.currentTarget.value})
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
                <input ref={priceMinRef} onChange={
                  (e) => {
                    setPriceInput({...priceInput, priceMin: +e.currentTarget.value})
                    setSearchQuery({...searchQuery, priceRange: priceInput})
                  }
                } type="range" min="0" max={MAX_PRICE_RANGE} value={priceInput.priceMin}/>
                {/* <div className="price_text">${priceInput.priceMin}</div> */}
                <input onChange={(e) => {
                  setPriceInput({...priceInput, priceMin: +e.currentTarget.value})
                }} className="price_text" type="text" value={priceInput.priceMin} />
              </div>
              <div>
                <label htmlFor="">max</label>
                <input ref={priceMaxRef} onChange={
                  (e) => {
                    setPriceInput({...priceInput, priceMax: +e.currentTarget.value})
                    setSearchQuery({...searchQuery, priceRange: priceInput})
                  }
                } type="range" min="0" max={MAX_PRICE_RANGE} value={priceInput.priceMax}/>
                {/* <div className="price_text">${priceInput.priceMax}</div> */}
                <input onChange={(e) => {
                  setPriceInput({...priceInput, priceMax: +e.currentTarget.value})
                }}  className="price_text" type="text" value={priceInput.priceMax} />
              </div>
            </div>
          </div>
        </div>

        <button className="clear-filter_btn" onClick={clearFilter}>Clear Filter</button>
      </div>

      <div className="search-bar">
        <label htmlFor="">Search for furniture:</label>
        <form onSubmit={(e) => onSearch(e)}>
          <input className="input_bar" onChange={
            (e) => setSearchQuery({...searchQuery, searchString: e.currentTarget.value})
          } type="text" />
          <input className="input_btn" type="submit" value="Search" />
        </form>
      </div>

    </div>
  )
}