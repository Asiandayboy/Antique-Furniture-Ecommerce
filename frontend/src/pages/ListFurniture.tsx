import { useState } from "react";
import Navbar from "../components/Navbar";


type FurnitureListing = {
  title: string,
  description: string,
  cost: number,
  type: string,
  style: string,
  condition: string,
  material: string,
}

const NONE: string = "None"
const EMPTY_LISTING: FurnitureListing = {
  title: "",
  description: "",
  cost: 0,
  type: NONE,
  material: NONE,
  style: NONE,
  condition: NONE
}


function validateForm(
  listing: FurnitureListing,
  imageFiles: File[]
): [err: boolean, keyName: string | null] {

  // validate form field text
  for (const key in listing) {
    const value = listing[key as keyof FurnitureListing]
    switch (typeof value) {
      case "string":
        if (value.trim() == "" || value.trim() == NONE) {
          return [true, key]
        }
        break
      case "number":
        if (value <= 0) {
          return [true, key]
        }
        break
    }
  }

  // validate form field images
  if (imageFiles.length == 0) {
    return [true, "images"]
  }

  return [false, null]
}


async function sendRequest(formData: FormData) {
  try {
    const res = await fetch("http://localhost:3000/list_furniture", {
      method: "POST",
      body: formData, // browser will automatically set the appropriate headers for the content
      credentials: "include"
    })

    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg)
    }

    const listingID = await res.text();
    console.log("The ID of your new listing is:", listingID)

  } catch(err) {
    console.error(err)
  }
}




export default function ListFurniture() {
  const [listing, setListing] = useState<FurnitureListing>(EMPTY_LISTING)
  const [imageFiles, setImageFiles] = useState<File[]>([])
  const [currentFormErr, setCurrentFormErr] = useState<string | null>()


  
  function onFormSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    const [ formErr, keyName ] = validateForm(listing, imageFiles)

    if (formErr) {
      const errStr = keyName![0].toUpperCase() + keyName?.slice(1) // convert keyname to uppercase
      setCurrentFormErr(errStr)
    } else {
      setCurrentFormErr(null)

      const formData = new FormData()
      formData.append("json_data", JSON.stringify(listing))
      imageFiles.forEach((file, i) => {
        formData.append("furniture_images", file)
      })

      sendRequest(formData)
    }

  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files) {
      setImageFiles(Array.from(e.target.files))
    }
  }

  
  return (
    <>
      <Navbar />
      <main>
        <div className="list-furniture_wrapper">
          <h1>Create a furniture listing</h1>
          {currentFormErr && 
            <div className="list-form-err">
              Form Error: Furniture {currentFormErr} is missing
            </div>
          }
          <form className="list-furniture-form" onSubmit={onFormSubmit} encType="multipart/form-data">
            <div className="list-title">
              <label htmlFor="title">Title</label>
              <input 
                type="text" 
                name="title"
                title="Title of the furniture listing"
                placeholder="Provide a title of your furniture listing"
                value={listing.title}
                onChange={(e) => {
                  setListing({...listing, title: e.currentTarget.value})
                }}
              />
            </div>
            <div className="list-description">
              <label htmlFor="description">Description</label>
              <textarea 
                name="description" 
                id="description" 
                placeholder="Provide a description of your furniture"
                value={listing.description}
                onChange={(e) => {
                  setListing({...listing, description: e.currentTarget.value})
                }}
              />
            </div>
            <div className="list-cost">
              <label htmlFor="cost">Cost</label>
              <input 
                type="number"
                name="cost"
                placeholder="Provide a cost for your furniture"
                id="cost" 
                value={listing.cost !== 0 ? listing.cost : ''}
                onChange={(e) => {
                  setListing({...listing, cost: +e.currentTarget.value})
                }}
              />
            </div>
            <div className="list-metadata">
              <div className="list-type">
                <label htmlFor="type">Furniture Type</label>
                <select 
                  name="type" 
                  id="type" 
                  value={listing.type}
                  onChange={(e) => {
                    setListing({...listing, type: e.currentTarget.value})
                  }}
                >
                  <option value={NONE}>{NONE}</option>
                  <option value="Bed">Bed</option>
                  <option value="Table">Table</option>
                  <option value="Nightstand">Nightstand</option>
                  <option value="Chair">Chair</option>
                  <option value="Cabinet">Cabinet</option>
                  <option value="Chest">Chest</option>
                  <option value="Armoire">Armoire</option>
                  <option value="Misc">Miscellaneous</option>
                </select>
              </div>
              <div className="list-material">
                <label htmlFor="material">Furniture Material</label>
                <select name="material" id="material" onChange={(e) => {
                    setListing({...listing, material: e.currentTarget.value})
                  }}>
                  <option value={NONE}>{NONE}</option>
                  <option value="Tiger Maple">Tiger Maple</option>
                  <option value="Maple">Maple</option>
                  <option value="Oak">Oak</option>
                  <option value="Mahogany">Mahogany</option>
                  <option value="Chestnut">Chestnut</option>
                  <option value="Pine">Pine</option>
                  <option value="Rosewood">Rosewood</option>
                  <option value="Cherry">Cherry</option>
                  <option value="Birch">Birch</option>
                  <option value="Walnut">Walnut</option>
                  <option value="Other">Other</option>
                </select>
              </div>
              <div className="list-style">
                <label htmlFor="style">Furniture Style</label>
                <select name="style" id="style" onChange={(e) => {
                    setListing({...listing, style: e.currentTarget.value})
                  }}>
                  <option value={NONE}>{NONE}</option>
                  <option value="Sheraton">Sheraton</option>
                  <option value="English">English</option>
                  <option value="Victorian">Victorian</option>
                  <option value="Baroque">Baroque</option>
                  <option value="Federal">Federal</option>
                  <option value="Rococo">Rococo</option>
                  <option value="Farmhouse">Farmhouse</option>
                  <option value="Contemporary">Contemporary</option>
                  <option value="Other">Other</option>
                </select>
              </div>
              <div className="list-condition">
                <label htmlFor="conditon">Furniture Condition</label>
                <select name="condition" id="condition" onChange={(e) => {
                    setListing({...listing, condition: e.currentTarget.value})
                  }}>
                  <option value={NONE}>{NONE}</option>
                  <option value="Mint">Mint</option>
                  <option value="Excellent">Excellent</option>
                  <option value="Good">Good</option>
                  <option value="Aged">Aged</option>
                  <option value="Restored">Restored</option>
                  <option value="Original Finish">Original Finish</option>
                </select>
              </div>
            </div>
            <div className="list-images">
              <label htmlFor="">Images</label>
              <div className="image-info">
                <label className="image-input" htmlFor="images">Choose Image Files</label>
                <input 
                  type="file"
                  name="images"
                  id="images" 
                  accept="image/jpeg,image/png"
                  onChange={handleFileChange}
                  multiple
                />
                <div className="image-files_wrapper">
                  <div className="files-selected_header">Files selected:</div>
                  <ul className="files-list">
                    {imageFiles.map((file, i) => (
                      <li key={i}>
                        {file.name}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
            <button className="list-finish_btn" type="submit">Finish Furniture Listing</button>
          </form>
        </div>
      </main>
    </>
  )
}