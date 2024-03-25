import { useState } from "react"

type Props = {
  imageURLs: string[]
}

export default function ImageSlider({ imageURLs }: Props) {
  const [currentImageIdx, setCurrentImageIdx] = useState<number>(0);

  return (
    <div className="image-slider_wrapper">
      <div className="slider-image">
        <div className="img_wrapper">
          <div className="img-hover">Click to expand</div>
          <img src={imageURLs[currentImageIdx]} alt={`slider image ${currentImageIdx+1}`} />
        </div>
      </div>
      <div className="slider-mini">
        {imageURLs.map((URL, i) => (
          <div 
            onClick={() => setCurrentImageIdx(i)} 
            className="img_wrapper"
            style={currentImageIdx == i && {
              "opacity": 1,
              "boxShadow": "0 5px 7px rgb(147, 147, 147)",
              "transform": "translateY(-5px)"
            } || {}}
          >
            <img src={URL} alt={`slider image ${i}`} />
          </div>
        ))}
      </div>
    </div>
  )
}