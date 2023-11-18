
import sisi from "k6/x/sisi";


const url     = "http://localhost:5120/api/simpleslideinterface/v1/slide/open/local/%s"
const rootPath = "C:\\Program Files\\3DHISTECH\\SimpleSlideInterface\\Examples\\DemoSlides\\%s"
const path    = "3DHISTECH"

export default function () {
    const token = sisi.getSlideToken(url,rootPath,path)
    console.log(token)
    const props = sisi.getBasicProperties(token)
    console.log(props)
    sisi.getTile(token,20,20,0)
}
