# Iolaus

> On Lerna's murderous hound, the many-headed hydra, he set his
> branding-iron...
>
> *Herakles*, Euripides

What is Iolaus?
---------------

If you're considering installing iolaus, you should first be aware
that iolaus is *completely unstable* and *completely unsupported*.
You are considering using experimental software that its author does
not consider to be fit for general consumption (and has therefore
never publicly announced) and has no time and little desire to make
fit for your consumption.  Not only that, but its build process
depends on software that is also unstable, which is also is not yet
intended for general use.

Iolaus is a [git](http://git-scm.org) porcelain which uses some ideas
from an older SCM I created called [darcs](http://darcs.net).  Thus
Iolaus is kin to Easy Git, Yap, darcs-git, vng, StGIT, Guilt and
TopGit.  Of these, I believe that Iolaus is closest in design, aims
and philosophy to [TopGit](http://repo.or.cz/w/topgit.git).  Clearly
there is an abundance of porcelains, and it's doubtful whether the
world needs yet another porcelain, particularly as there are already
two that advertise a darcs-like interface.  On the other hand, the
abundance of porcelains for git reflect how easy they are to create,
and I am accustomed to using a SCM that I'm very familiar with.

The command set of Iolaus is reasonably similar to that of darcs, but
I plan to gradually work it into better agreement with git's flags and
terminology.

Iolaus is still very experimental, and I don't recommend it for anyone
who isn't interested in hacking on Iolaus itself.

Who is Iolaus?
--------------

Iolaus is the nephew of Heracles (also known by his latin name
Hercules), who went on adventures with him and helped with some of his
labors.  Most famously, when Heracles was fighting the hydra, he had
trouble because every time he knocked one head off, two grew in his
place.  He called to Iolaus for help, and Iolaus figured out that if
he cauterized each neck with his torch after Heracles knocked the head
off, they would stop sprouting.  There was still the problem of the
hydra's one immortal head, but they just hid that under a rock.

> You are running away, Socrates, said Dionysodorus, and refusing to
> answer.
>
> No wonder, I said, for I am not a match for one of you, and a
> fortiori I must run away from two. I am no Heracles; and even
> Heracles could not fight against the Hydra, who was a she-Sophist,
> and had the wit to shoot up many new heads when one of them was cut
> off; especially when he saw a second monster of a sea-crab, who was
> also a Sophist, and appeared to have newly arrived from a
> sea-voyage, bearing down upon him from the left, opening his mouth
> and biting. When the monster was growing troublesome he called
> Iolaus, his nephew, to his help, who ably succoured him; but if my
> Iolaus, who is my brother Patrocles [the statuary], were to come, he
> would only make a bad business worse.
>
> *Euthydemus*, Plato

The origins of Iolaus
---------------------

I realized that the semantics of git are actually not nearly so far
from those of darcs as I had previously thought.  In particular, if we
view each commit as describing a patch in its "primitive context" (to
use darcs-speak), then there is basically a one-to-one mapping from
darcs' semantics to a git repository.   The catch is that it must be a 
git repository with *multiple heads*!

Fortunately, this is not such a foreign concept to git.  In fact, git
has a whole framework to help users manage repositories with multiple
heads (see, e.g. checkout and branch).  So it's not so very foreign at
all.  There are just a couple of major differences how git works.
First, in git your working directory will only reflect *one* of the
heads, while in darcs (or iolaus) the working directory reflects the
union of all changes in the repository.

The other major difference is that in git, new commits will normally
not introduce new heads, meaning that the history is normally linear.
This inhibits "cherry picking", since you cannot pull a patch without
pulling its parent.  Iolaus instead commutes patches back, so that the
parents of each commit are only the patches that that patch depends
upon.  In practice, we don't want to commute *that* far back, so we
only commute a far back as we might want to cherry pick.

Legendary head management
-------------------------

> ...and in his turn called for help on Iolaus who, by setting fire to
> a piece of the neighbouring wood and burning the roots of the heads
> with the brands, prevented them from sprouting. Having thus got the
> better of the sprouting heads, he chopped off the immortal head, and
> buried it, and put a heavy rock on it, beside the road that leads
> through Lerna to Elaeus.
>
> *Apollodorus*, translated by Sir James G. Frazer, 1921.

Iolaus takes charge of creating and decapitating heads.  This makes
its name something of a misnomer, since Heracles was in charge of
removing heads, which grew spontaneously, and Iolaus was merely tasked
with preventing their regrowth.  Nevertheless, iolaus the program
takes charge of both the growth and decapitation of heads, allowing
you (Heracles?) the programmer to focus your attention on
code---that's the theory, anyhow.
