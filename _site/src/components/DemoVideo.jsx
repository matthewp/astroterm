import mp4Url from '../assets/videos/demo.mp4?url';
import webmUrl from '../assets/videos/demo.webm?url';
import '../styles/DemoVideo.css';

const dev = import.meta.env.DEV
const base = dev ? '' : '/astroterm'

export default function() {
    return (
        <video class="demo-video" controls>
            <source src={`${base}${mp4Url}`} type="video/mp4"/>
            <source src={`${base}${webmUrl}`} type="video/webm"/>
        </video>
    );
}