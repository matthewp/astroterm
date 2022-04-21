import mp4Url from '../assets/videos/demo.mp4?url';
import webmUrl from '../assets/videos/demo.webm?url';
import '../styles/DemoVideo.css';

export default function() {
    return (
        <video class="demo-video" controls>
            <source src={`/astroterm${mp4Url}`} type="video/mp4"/>
            <source src={`/astroterm${webmUrl}`} type="video/webm"/>
        </video>
    );
}