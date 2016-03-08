<?php


// ** MySQL settings ** //
/** The name of the database for WordPress */
define('DB_NAME', 'wordpress');

/** MySQL database username */
define('DB_USER', 'wordpress');

/** MySQL database password */
define('DB_PASSWORD', 'wordpress');

/** MySQL hostname */
define('DB_HOST', 'localhost');

/** Database Charset to use in creating database tables. */
define('DB_CHARSET', 'utf8');

/** The Database Collate type. Don't change this if in doubt. */
define('DB_COLLATE', '');

define('AUTH_KEY',         '}Fg)KrSu{9$p`OF;GUPub9~Y-|$R?a|,EwT{!k=qdc,=6Zn^.O?^W[P(XrX1L@2R');
define('SECURE_AUTH_KEY',  '>gnr}q?E(WiLX0<gh*EM]q7GGSwb#4,0PHD8)Zb;aXZ9=$<d9=d|sJ+!nrQ6@cI.');
define('LOGGED_IN_KEY',    '4[jbIqjxTQT$s)|#H9s!jJgQ/[Fi;AHQq.K2)]njFY+;tm8$^C|xcm!I&-)Ot+5B');
define('NONCE_KEY',        'J{Nbv7MKWb9U;d?w8P:[6p,)W$yoXgh!&@e#bq *+!n>no-!Nm.Ry0~#BkWoPLYX');
define('AUTH_SALT',        '-XR0pWI&HHaPE|iQEcDKLXf4u`|Rj.sj(tc%t;>909$^~3[~@fzm+u6&`(U[`B!|');
define('SECURE_AUTH_SALT', 'V0/CM[pW*#++]Hblx9Y_(ZV&zctv(hbBee[.<RWOAQe=1#E0#9WIf{LPs2_5PVU7');
define('LOGGED_IN_SALT',   'f).CF/UQGQDR`aE6c(Uwh1iF|K`jv/-U{,vGQCNp .R)Xc|3f40gs+z.9*1q_fmg');
define('NONCE_SALT',       'L=KvSc+&j.HMu}7<Xn&eiG-eyzNT6+Ii_2+ A[oPl3qIrtrs;93CofOv|>e TcKe');


$table_prefix = 'wp_';


define( 'WP_HOME', 'http://vccw.dev' );
define( 'WP_SITEURL', 'http://vccw.dev' );
define( 'JETPACK_DEV_DEBUG', true );
define( 'WP_DEBUG', true );
define( 'FORCE_SSL_ADMIN', false );
define( 'SAVEQUERIES', false );



/* That's all, stop editing! Happy blogging. */

/** Absolute path to the WordPress directory. */
if ( !defined('ABSPATH') )
	define('ABSPATH', dirname(__FILE__) . '/');

/** Sets up WordPress vars and included files. */
require_once(ABSPATH . 'wp-settings.php');
