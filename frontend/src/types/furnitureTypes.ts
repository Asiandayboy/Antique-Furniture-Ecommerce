export type FurnitureType = "bed"
| "table"
| "desk"
| "chair"
| "chest"
| "nightstand"
| "cabinet"
| "lamp"

export type FurnitureStyle = "english"
| "sheraton"
| "victorian"
| "baroque"
| "federal"
| "sheraton"

export type FurnitureMaterial = "tiger maple"
| "cherry"
| "walnut"
| "mahogany"
| "oak"
| "maple"
| "chestnut"
| "pine"
| "rosewood"
| "birch"

export type FurnitureBedSize = "california king"
| "king" 
| "queen"
| "double"
| "twin"

/**Returns true if the string is a furniture type, else false */
export function isFurnitureType(str: string): str is FurnitureType {
  return ['bed', 'table', 'desk', 'chair', 'chest', 'nightstand', 'cabinet', 'lamp'].includes(str);
}

/**Returns true if the string is a furniture style, else false */
export function isFurnitureStyle(str: string): str is FurnitureStyle {
  return ['english', 'sheraton', 'victorian', 'baroque', 'federal', 'sheraton'].includes(str);
}

/**Returns true if the string is a furniture material, else false */
export function isFurnitureMaterial(str: string): str is FurnitureMaterial {
  return ['tiger maple', 'cherry', 'walnut', 'mahogany', 'oak', 'maple', 'chestnut', 'pine', 'rosewood', 'birch'].includes(str);
}

/**Returns true if the string is a type of bed size, else false */
export function isFurnitureBedSize(str: string): str is FurnitureBedSize {
  return ['california king', 'king', 'queen', 'double', 'twin'].includes(str);
}