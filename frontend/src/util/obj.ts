/**Compares objects by content instead of by reference */
export function deepEqual(obj1: any, obj2: any) {
      // If both objects are not objects, they can't be equal
  if (typeof obj1 !== 'object' || typeof obj2 !== 'object') {
    return false;
  }

  // If number of keys in obj1 and obj2 are different, they can't be equal
  if (Object.keys(obj1).length !== Object.keys(obj2).length) {
    return false;
  }

  // Check if all properties of obj1 are present in obj2 with equal values
  for (let key in obj1) {
    if (!obj2.hasOwnProperty(key) || obj1[key] !== obj2[key]) {
      return false;
    }
  }

  // Check if all properties of obj2 are present in obj1 with equal values
  for (let key in obj2) {
    if (!obj1.hasOwnProperty(key) || obj2[key] !== obj1[key]) {
      return false;
    }
  }

  // If all checks passed, the objects are equal
  return true;
}